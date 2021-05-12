package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestReadConf(t *testing.T) {
	t.Parallel()
	// build our expected data
	var nur newUserRequest
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

func Example_validateFiles() {
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

func ExampleValidate() {
	s := serviceSpec{APIVersion: "openslo/v1alpha"}
	validate(s)

	// Output:
	// Valid
}
