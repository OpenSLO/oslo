package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var data = `
conf:
  username: slobear
  name: Oslo Joe
  age: 21
  password: superstrong
`

func TestReadConf(t *testing.T) {
	// build our expected data
	nur := &NewUserRequest{}
	err := yaml.Unmarshal([]byte(data), nur)

	// make sure we don't get any errors
	assert.Nil(t, err)

	// read the file
	c, _ := readConf("../test/valid.yaml")

	assert.NotNil(t, c)
	assert.Equal(t, nur, c)
}

func ExampleValidate() {
	a := []string{"../test/valid.yaml"}
	validateFiles(a)

	b := []string{"../test/invalid.yaml"}
	validateFiles(b)

	c := []string{"../test/missing.yaml"}
	validateFiles(c)

	// Unordered output:
	// Valid!
	// Invalid
	//   - Conf.Age (less than min)
	// Invalid
	//   - Conf.Username (zero value, less than min)
	//   - Conf.Name (zero value)
	//   - Conf.Password (zero value, less than min)

}
