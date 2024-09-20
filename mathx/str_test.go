package mathx

import (
	"fmt"
	"testing"

	"github.com/gophero/goal/assert"
)

func TestTrim0(t *testing.T) {
	a := Ceilrf(3.14, 5)
	fmt.Println(Trim0(a))
	assert.True(Trim0(a) == "3.14")
	a = "1.0"
	assert.True(Trim0(a) == "1")
	a = "0.1"
	assert.True(Trim0(a) == a)
	a = ""
	assert.True(Trim0(a) == a)
	a = "100"
	assert.True(Trim0(a) == a)
	a = "0.000001"
	assert.True(Trim0(a) == a)
}
