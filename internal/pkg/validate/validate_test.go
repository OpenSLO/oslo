/*
Copyright © 2021 OpenSLO Team

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
package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_readConf(t *testing.T) {
	t.Parallel()

	c, e := ReadConf("../../../test/valid-service.yaml")

	assert.NotNil(t, c)
	assert.Nil(t, e)

	_, e = ReadConf("../../../test/non-existent.yaml")

	assert.NotNil(t, e)
}

func Test_validateFiles(t *testing.T) {
	t.Parallel()

	validFiles := []struct {
		filename string
	}{
		{"../../../test/valid-service.yaml"},
		{"../../../test/valid-slos-ratio.yaml"},
		{"../../../test/valid-slos-threshold.yaml"},
	}

	for _, tt := range validFiles {
		tt := tt
		t.Run(tt.filename, func(t *testing.T) {
			t.Parallel()
			a := []string{tt.filename}
			assert.Nil(t, validateFiles(a))
		})
	}

	d := []string{"../../test/invalid-service.yaml"}
	assert.NotNil(t, validateFiles(d))
}
