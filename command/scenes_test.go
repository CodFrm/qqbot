package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScenesList(t *testing.T) {
	_, err := ScenesList("", 1)
	assert.Nil(t, err)
}
