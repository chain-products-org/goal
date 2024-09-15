package iox_test

import (
	"fmt"
	"github.com/gophero/goal/iox"
	"strings"
	"testing"
)

func TestWalkAllFiles(t *testing.T) {
	fs := iox.WalkDir("/Users/home/Downloads")
	fmt.Println(strings.Join(fs, "\n"))
}
