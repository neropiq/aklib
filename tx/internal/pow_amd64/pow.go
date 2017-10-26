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

#define SHA256_WORDS 8
#define SHA256_DIGEST_SIZE 32
#define NonceLocation 4

#ifndef _MSC_VER
#include <pthread.h>
#endif

#ifdef _WIN32
HANDLE	mutex = CreateMutex(NULL,FALSE,NULL);
#else
pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;
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

int pwork(uint8_t dat[], uint64_t inlen, uint8_t difficulty)
{
    int i;
    uint8_t hash[32], hash2[32];

    while (1) {
    loop:
        blake2b(hash, 32, dat, inlen, NULL, 0);
        sha256_32bytes(hash, hash2);
        if (isValid(hash2, difficulty)) {
            return 1;
        }
        for (i = 0; i < 32; i++) {
            dat[NonceLocation + i]++;
            if (dat[NonceLocation + i] != 0x00) {
                goto loop;
                ;
            }
        }

        return 0;
    }
}

*/
import "C"
import "unsafe"

//PoW does PoW for amd64 platforms..
func PoW(dat []byte, difficulty byte) bool {
	r := C.pwork((*C.uint8_t)(unsafe.Pointer(&dat[0])), C.uint64_t(len(dat)), C.uint8_t(difficulty))
	return r != 0
}
