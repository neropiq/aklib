// +build linux,amd64

// Copyright (c) 2017 Aidos Developer

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package pow64

// #cgo CFLAGS: -O3  -Wall -Wextra -pedantic -Wno-long-long
/*
#include <string.h>
#include <stdint.h>
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>

#define SHA256_WORDS 8
#define SHA256_DIGEST_SIZE 32
#define NonceLocation 4

#ifndef _MSC_VER
#include <pthread.h>
#endif

#ifdef _WIN32
#define MUTEX HANDLE
#define NEW_MUTEX CreateMutex(NULL, FALSE, NULL)
#define LOCK(mutex) WaitForSingleObject(mutex, INFINITE)
#define UNLOCK(mutex) ReleaseMutex(mutex)
#define PTHREAD_T HANDLE
#else
#define MUTEX pthread_mutex_t
#define NEW_MUTEX PTHREAD_MUTEX_INITIALIZER;
#define LOCK(mutex) pthread_mutex_lock(&mutex)
#define UNLOCK(mutex) pthread_mutex_unlock(&mutex)
#define PTHREAD_T pthread_t
#endif

extern void sha256_sse4(void* input_data, uint32_t digest[8], uint64_t num_blks);
extern int blake2b(void* out, size_t outlen, const void* in, size_t inlen, const void* key, size_t keylen);

inline void sha256_32bytes(const void* data, uint8_t hash[SHA256_DIGEST_SIZE])
{
    int i;
    uint32_t buf[16] = { 0 };
    buf[8] = 0x80;
    buf[15] = 0x10000;

    uint32_t schash[] = {
        0x6a09e667,
        0xbb67ae85,
        0x3c6ef372,
        0xa54ff53a,
        0x510e527f,
        0x9b05688c,
        0x1f83d9ab,
        0x5be0cd19
    };

    memcpy(buf, data, 32);
    sha256_sse4(buf, schash, 1);

    for (i = 0; i < SHA256_WORDS; i++) {
        hash[0] = (uint8_t)(schash[i] >> 24);
        hash[1] = (uint8_t)(schash[i] >> 16);
        hash[2] = (uint8_t)(schash[i] >> 8);
        hash[3] = (uint8_t)schash[i];
        hash += 4;
    }
}

inline int isValid(uint8_t hash[], uint8_t difficulty)
{
    uint8_t i, d, b;

    for (i = 0; i<difficulty>> 3; i++) {
        if (hash[31 - i] != 0x00) {
            return 0;
        }
    }
    d = difficulty - (i << 3);
    if (d == 0) {
        return 1;
    }
    b = (1 << (8 - d)) - 1;
    if (hash[31 - i] > b) {
        return 0;
    }
    return 1;
}

void add(uint8_t* b, int inc)
{
    int i, j = 0;
    for (j = 0; j < inc; j++) {
        for (i = 0; i < 32; i++) {
            b[NonceLocation + 31 - i]++;
            if (b[NonceLocation + 31 - i] != 0x00) {
                return;
            }
        }
    }
}

typedef struct param {
    uint8_t* dat;
    uint64_t inlen;
    uint64_t difficulty;
    int n;
    MUTEX m;
    int* stop;
} PARAM;

#ifdef _WIN32
unsigned __stdcall pwork_(void* p)
#else
void* pwork_(void* p)
#endif
{
    PARAM* par = (PARAM*)(p);
    int i;
    uint8_t hash[32], hash2[32];
    while (1) {
    loop:
        LOCK(par->m);
        int stop_ = *(par->stop);
        UNLOCK(par->m);
        if (stop_) {
            break;
        }
        blake2b(hash, 32, par->dat, par->inlen, NULL, 0);
        sha256_32bytes(hash, hash2);
        if (isValid(hash2, par->difficulty)) {
            LOCK(par->m);
            *(par->stop) = par->n;
            UNLOCK(par->m);
            break;
        }
        for (i = 0; i < 32; i++) {
            par->dat[NonceLocation + i]++;
            if (par->dat[NonceLocation + i] != 0x00) {
                goto loop;
            }
        }
        break;
    }
#ifdef _WIN32
    return 0;
#else
    return NULL;
#endif
}

int getCpuNum()
{
#if defined(__linux) || defined(__APPLE__)
    // for linux
    return sysconf(_SC_NPROCESSORS_ONLN);
#elif defined(__MINGW64__) || defined(__MINGW32__) || defined(_MSC_VER)
    // for windows and wine
    SYSTEM_INFO info;
    GetSystemInfo(&info);
    return info.dwNumberOfProcessors;
#endif
}

int pwork(uint8_t dat[], uint64_t inlen, uint8_t difficulty)
{
    int i = 0, stop = 0;
    int procs = getCpuNum();

    PTHREAD_T* thread = (PTHREAD_T*)calloc(sizeof(PTHREAD_T), procs);
    MUTEX mut = NEW_MUTEX;

    PARAM* p = (PARAM*)calloc(sizeof(PARAM), procs);
    for (i = 0; i < procs; i++) {
        p[i].dat = (uint8_t*)malloc(inlen);
        memcpy(p[i].dat, dat, inlen);
        add(p[i].dat, i);
        p[i].inlen = inlen;
        p[i].difficulty = difficulty;
        p[i].n = i + 1;
        p[i].stop = &stop;
        p[i].m = mut;
#ifdef _WIN32
        unsigned int id = 0;
        thread[i] = (HANDLE)_beginthreadex(NULL, 0, pwork_, (LPVOID)&p[i], 0, NULL);
        if (thread[i] == NULL) {
#else
        int ret = pthread_create(&thread[i], NULL, pwork_, &p[i]);
        if (ret != 0) {
#endif
            fprintf(stderr, "can not create thread\n");
            return 1;
        }
    }
    for (i = 0; i < procs; i++) {
#ifdef _WIN32
        int ret = WaitForSingleObject(thread[i], INFINITE);
        CloseHandle(thread[i]);
        if (ret == WAIT_FAILED) {
#else
        int ret = pthread_join(thread[i], NULL);
        if (ret != 0) {
#endif
            fprintf(stderr, "can not join thread\n");
            return 1;
        }
    }
    if (stop) {
        memcpy(&dat[NonceLocation], &(p[stop - 1].dat[NonceLocation]), 32);
    }
    for (i = 0; i < procs; i++) {
        free(p[i].dat);
    }
    free(thread);
    free(p);
    if (!stop) {
        return 0;
    }
    return 1;
}
*/
import "C"
import "unsafe"

//PoW does PoW for amd64 platforms..
func PoW(dat []byte, difficulty byte) bool {
	r := C.pwork((*C.uint8_t)(unsafe.Pointer(&dat[0])), C.uint64_t(len(dat)), C.uint8_t(difficulty))
	return r != 0
}
