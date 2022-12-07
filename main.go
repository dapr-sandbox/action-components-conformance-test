/*
Copyright 2022 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	dir := flag.String("dir", ".", "directory")
	flag.Parse()

	cmd := exec.Command("go", "test", "-run", "TestPluggableConformance/"+os.Getenv("INPUT_TYPE"), "-v")

	stdout, err := cmd.StdoutPipe()
	cmd.Dir = *dir
	cmd.Env = os.Environ()
	cmd.Stderr = cmd.Stdout

	if err != nil {
		panic(err)
	}
	if err = cmd.Start(); err != nil {
		panic(err)
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Println(string(tmp)) //nolint:forbidigo
		if err != nil {
			break
		}
	}

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}
