package command

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestIsWordImage(t *testing.T) {
	b, _ := ioutil.ReadFile("./data/1.jpg")
	s, err := IsWordImage(0, 0, b)
	assert.Nil(t, err)
	assert.Equal(t, "", s)

	b, _ = ioutil.ReadFile("./data/2.jpg")
	s, err = IsWordImage(0, 0, b)
	assert.Nil(t, err)
	assert.Equal(t, "OK", s)

	b, _ = ioutil.ReadFile("./data/3.jpg")
	s, err = IsWordImage(0, 0, b)
	assert.Nil(t, err)
	assert.Equal(t, "OK", s)

	b, _ = ioutil.ReadFile("./data/4.jpg")
	s, err = IsWordImage(0, 0, b)
	assert.Nil(t, err)
	assert.Equal(t, "OK", s)

	b, _ = ioutil.ReadFile("./data/5.jpg")
	s, err = IsWordImage(0, 0, b)
	assert.Nil(t, err)
	assert.Equal(t, "OK", s)

	b, _ = ioutil.ReadFile("./data/6.jpg")
	s, err = IsWordImage(0, 0, b)
	assert.Nil(t, err)
	assert.Equal(t, "", s)
}
