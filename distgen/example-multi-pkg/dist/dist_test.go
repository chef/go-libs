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

package dist_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	subject "github.com/chef/go-libs/distgen/example-multi-pkg/dist"
)

func TestDistDefaults(t *testing.T) {
	assert.Equal(t,
		"dist-test", subject.Exec,
		"something is going on with the local dist variables",
	)
	assert.Equal(t,
		"bar", subject.Foo,
		"something is going on with the generation of distributable variable, check the 'distgen' lib",
	)
}
