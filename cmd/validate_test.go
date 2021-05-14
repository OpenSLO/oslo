package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestReadConf(t *testing.T) {
	t.Skip()
	t.Parallel()
	// build our expected data
	var nur interface{}
	data := []byte(`
conf:
  username: slobear
  name: Oslo Joe
  age: 21
  password: superstrong
`)
	err := yaml.Unmarshal(data, &nur)

	assert.Nil(t, err)

	c, _ := readConf("../test/valid.yaml")

	assert.NotNil(t, c)
	assert.Equal(t, nur, c)
}

func Test_validateFiles(t *testing.T) {
	a := []string{"../test/valid-service.yaml"}
	assert.Nil(t, validateFiles(a))
	b := []string{"../test/valid-slos-ratio.yaml"}
	assert.Nil(t, validateFiles(b))
	c := []string{"../test/valid-slos-threshold.yaml"}
	assert.Nil(t, validateFiles(c))
	d := []string{"../test/invalid-service.yaml"}
	assert.NotNil(t, validateFiles(d))

	// b := []string{"../test/invalid.yaml"}
	// validateFiles(b)

	// c := []string{"../test/missing.yaml"}
	// validateFiles(c)

	// Unordered output:
	// Valid!
	// Invalid
	//   - Conf.Age (less than min)
	// Invalid
	//   - Conf.Username (zero value, less than min)
	//   - Conf.Name (zero value)
	//   - Conf.Password (zero value, less than min)
}
