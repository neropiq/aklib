// Copyright (c) 2018 Aidos Developer

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

package updater

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

//Update start to update the software.
func Update(repo, version string) error {
	latest, found, err := selfupdate.DetectLatest(repo)
	if err != nil {
		return err
	}

	v := semver.MustParse(version)
	if !found || latest.Version.Equals(v) {
		log.Println("Current version is the latest")
		return nil
	}

	for ok := false; !ok; {
		fmt.Print("Do you want to update to", latest.Version, "? (y/n): ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return err
		}
		switch input {
		case "y\n":
			ok = true
		case "n\n":
			return nil
		default:
			fmt.Println("Invalid input")
		}
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, os.Args[0]); err != nil {
		return err
	}
	log.Println("Successfully updated to version", latest.Version)
	return nil
}
