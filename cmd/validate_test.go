package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_readConf(t *testing.T) {
	t.Parallel()

	c, e := readConf("../test/v1alpha1_valid-service.yaml")

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
		{"../test/v1alpha1_valid-service.yaml"},
		{"../test/v1alpha1_valid-slos-ratio.yaml"},
		{"../test/v1alpha1_valid-slos-threshold.yaml"},
		{"../test/v1beta1_valid-slos-ratio.yaml"},
		{"../test/v1beta1_valid-sli-ratio.yaml"},
	}

	for _, tt := range validFiles {
		tt := tt
		t.Run(tt.filename, func(t *testing.T) {
			t.Parallel()
			a := []string{tt.filename}
			assert.Nil(t, validateFiles(a))
		})
	}

	d := []string{"../test/v1alpha1_invalid-service.yaml"}
	assert.NotNil(t, validateFiles(d))
}
