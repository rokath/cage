package cage

import (
	"fmt"
	"os"
	"testing"

	"github.com/udhos/equalfile"
)

func equalFileContent(fn0, fn1 string) bool {
	cmp := equalfile.New(nil, equalfile.Options{}) // compare using single mode
	ok, err := cmp.CompareFile(fn0, fn1)
	if nil != err {
		ok = false
	}
	return ok
}

func equalFiles(t *testing.T, fn0, fn1 string) {
	ok := equalFileContent(fn0, fn1)
	if false == ok {
		t.FailNow()
	}
}
func TestStart(t *testing.T) {
	lfn := "testdata/logfile.txt"
	os.Remove(lfn)

	c := Start(lfn)
	fmt.Println("test0")
	fmt.Fprintln(os.Stderr, "test1")
	Stop(c)

	equalFiles(t, lfn, "testdata/test01.txt")

	c = Start(lfn)
	fmt.Println("test2")
	fmt.Println("test3")
	Stop(c)

	equalFiles(t, lfn, "testdata/test0123.txt")
}
