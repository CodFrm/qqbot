package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCloudflareResolve(t *testing.T) {
	ip, err := CloudflareResolve("blog.icodef.com")
	assert.Nil(t, err)
	println(ip)
}
