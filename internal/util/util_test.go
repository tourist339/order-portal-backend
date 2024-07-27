package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGenerateUniqueID(t *testing.T) {
	t.Run("check prefix", func(t *testing.T) {
		pre := "UT"
		s := GenerateUniqueID(pre)
		fmt.Println(s)
		assert.True(t, strings.HasPrefix(s, pre))
	})
}
