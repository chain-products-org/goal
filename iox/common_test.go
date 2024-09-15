package iox_test

import (
	"github.com/gophero/goal/iox"
	"github.com/gophero/goal/testx"
	"testing"
)

func TestExists(t *testing.T) {
	lg := testx.Wrap(t)

	lg.Case("give a existing dir")
	f := ".."
	lg.Require(iox.Exists(f), "should exist")

	lg.Case("give a existing file")
	f = "./file_test.go"
	lg.Require(iox.Exists(f), "should exist")
}

func TestIsDir(t *testing.T) {
	lg := testx.Wrap(t)

	lg.Case("give a existing dir")
	f := "../iox"
	lg.Require(iox.IsDir(f), "is dir")

	lg.Case("give a existing file")
	f = "../iox/file_test.go"
	lg.Require(!iox.IsDir(f), "is not a dir")
}

func TestIsRegularFile(t *testing.T) {
	lg := testx.Wrap(t)

	lg.Case("give a existing dir")
	f := "../iox/"
	lg.Require(!iox.IsRegularFile(f), "is not a regular file")

	lg.Case("give a existing regular file")
	f = "../iox/file_test.go"
	lg.Require(iox.IsRegularFile(f), "is a regular file")

	lg.Case("give a soft symlink file")
	f = "/etc"
	lg.Require(!iox.IsRegularFile(f), "is not a regular file")
}
