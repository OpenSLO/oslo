package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_readConf(t *testing.T) {
	t.Parallel()

	c, e := readConf("../test/valid-service.yaml")

	assert.NotNil(t, c)
	assert.Nil(t, e)

	_, e = readConf("../test/non-existent.yaml")

	assert.NotNil(t, e)
}

func Test_validateFiles(t *testing.T) {
	t.Parallel()

	validFiles := []struct {
		filename string
	}{
		{"../test/valid-service.yaml"},
		{"../test/valid-slos-ratio.yaml"},
		{"../test/valid-slos-threshold.yaml"},
	}

	for _, tt := range validFiles {
		tt := tt
		t.Run(tt.filename, func(t *testing.T) {
			t.Parallel()
			a := []string{tt.filename}
			assert.Nil(t, validateFiles(a))
		})
	}

	d := []string{"../test/invalid-service.yaml"}
	assert.NotNil(t, validateFiles(d))
}
