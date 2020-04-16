package cage

import (
	"fmt"
	"log"
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

	lfn := "./testdata/act.log"
	efn := "./testdata/exp.log"

	// clear
	os.Remove(lfn)
	os.Remove(efn)

	// write expectation
	efh, err := os.OpenFile(efn, os.O_RDWR|os.O_CREATE, 0666)
	ok(t, err)
	_, err = fmt.Fprintln(efh, "log00")
	ok(t, err)
	_, err = fmt.Fprintln(efh, "fmt01")
	ok(t, err)
	_, err = fmt.Fprintln(efh, "err02")
	ok(t, err)
	err = efh.Close()
	ok(t, err)

	// start logging
	c := Start(lfn)

	// produce logs
	log.SetFlags(0)
	fmt.Println("fmt01")
	fmt.Fprintln(os.Stderr, "err02")
	log.Println("log00")

	// expect written data before close
	equalFiles(t, lfn, efn)

	// stop logging
	Stop(c)

	// expect written data after close
	equalFiles(t, lfn, efn)

	// continue expectation
	efh, err = os.OpenFile(efn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	ok(t, err)
	_, err = fmt.Fprintln(efh, "log10")
	ok(t, err)
	_, err = fmt.Fprintln(efh, "fmt11")
	ok(t, err)
	_, err = fmt.Fprintln(efh, "err12")
	ok(t, err)
	err = efh.Close()
	ok(t, err)

	// continue logging
	c = Start(lfn)

	// produce logs
	log.SetFlags(0)
	fmt.Println("fmt11")
	fmt.Fprintln(os.Stderr, "err12")
	log.Println("log10")

	// expect written data before close
	equalFiles(t, lfn, efn)

	// stop logging
	Stop(c)

	// expect written data after close
	equalFiles(t, lfn, efn)

	// clean up
	err = os.Remove(lfn)
	ok(t, err)
	err = os.Remove(efn)
	ok(t, err)
}
