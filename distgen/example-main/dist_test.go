//
// Copyright 2020 Chef Software, Inc.
// Author: Salim Afiune <afiune@chef.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistDefaults(t *testing.T) {
	assert.Equal(t,
		"dist-test", Exec,
		"something is going on with the local dist variables",
	)
	assert.Equal(t,
		"local-var", LocalVar,
		"something is going on with the local dist variables",
	)
	assert.Equal(t,
		"chef-client", ClientExec,
		"something is going on with the generation of distributable variable, check the 'distgen' lib",
	)
}

func TestDoMain(t *testing.T) {
	// create a reader and writer to capture the output of the main() function
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	// substitute the STDOUT and restore the system STDOUT at the end of the test
	stdout := os.Stdout
	defer func() {
		os.Stdout = stdout
	}()
	os.Stdout = writer

	// create a channel to ready the buffer
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()

	// execute the main function
	main()

	writer.Close()

	assert.Equal(t,
		"dist-test:\n * Global product name: 'Chef Infra Server'\n * Local variable: 'local-var'\n",
		<-out,
		"something is going on with the local dist variables",
	)
}
