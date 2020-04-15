package cage

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/udhos/equalfile"
)

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

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
	lfn := "logfile.txt"
	efn := "explogf.txt"
	os.Remove(lfn)
	os.Remove(efn)

	efh, err := os.OpenFile(efn, os.O_RDWR|os.O_CREATE, 0666)
	ok(t, err)
	_, err = fmt.Fprintln(efh, "test0")
	ok(t, err)
	_, err = fmt.Fprintln(efh, "test1")
	ok(t, err)
	err = efh.Close()
	ok(t, err)

	c := Start(lfn)
	fmt.Println("test0")
	fmt.Fprintln(os.Stderr, "test1")
	Stop(c)

	equalFiles(t, lfn, efn)

	efh, err = os.OpenFile(efn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	ok(t, err)

	_, err = fmt.Fprintln(efh, "test2")
	ok(t, err)
	_, err = fmt.Fprintln(efh, "test3")
	ok(t, err)
	err = efh.Close()
	ok(t, err)

	c = Start(lfn)
	fmt.Println("test2")
	fmt.Println("test3")
	Stop(c)

	equalFiles(t, lfn, efn)

	err = os.Remove(lfn)
	ok(t, err)
	err = os.Remove(efn)
	ok(t, err)
}
