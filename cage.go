package cage

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"
)

// Container keeps re-direction informantion
type Container struct {
	oldLog io.Writer // old

	origStdout   *os.File // old
	writerStdout *os.File // new

	origStderr   *os.File // old
	writerStderr *os.File // new

	sChannel chan string // channel for os.Stdout messages
	eChannel chan string // channel for os.Stderr messages

	lfHandle *os.File // logfile handle
}

// Start does append all output parallel into a logfile with name fn
func Start(fn string) *Container {

	// start logging only if fn not "none"
	if "none" == fn {
		log.Println("No logfile")
		return nil
	}

	// open logfile
	lfH, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		//log.Fatalf("error opening file %s: %v", fn, err)
		fn = "off"
		return nil
	}

	// open pipes
	rStdout, wStdout, _ := os.Pipe()
	rStderr, wStderr, _ := os.Pipe()

	// create container for recovering
	c := &Container{
		oldLog: log.Writer(),

		origStdout:   os.Stdout,
		writerStdout: wStdout,

		origStderr:   os.Stderr,
		writerStderr: wStderr,

		sChannel: make(chan string),
		eChannel: make(chan string),

		lfHandle: lfH,
	}

	// re-direct
	log.SetOutput(io.MultiWriter(c.oldLog, c.lfHandle)) // writing to log will go also to logfile now

	os.Stdout = c.writerStdout // all to os.Stdout goees now to c.writerStdout and comee out of rStdout now
	os.Stderr = c.writerStderr // all to os.Stderr goees now to c.writerStderr and comes out of rStderr now

	// create duplication
	teeOut := io.MultiWriter(c.origStdout, c.lfHandle)
	teeErr := io.MultiWriter(c.origStderr, c.lfHandle)

	// pipes reader and channel writer
	go func(sout, eout chan string, readerStdout *os.File, readerStderr *os.File) {
		var bufStdout bytes.Buffer
		// _, _ = io.Copy(teeOut, readerStdout) // CHECK this
		_, _ = io.Copy(&bufStdout, readerStdout)
		if bufStdout.Len() > 0 {
			sout <- bufStdout.String()
		}

		var bufStderr bytes.Buffer
		_, _ = io.Copy(&bufStderr, readerStderr)
		if bufStderr.Len() > 0 {
			eout <- bufStderr.String()
		}
	}(c.sChannel, c.eChannel, rStdout, rStderr)

	// channel reader and multi writer
	go func(c *Container) {
		for {
			select {
			case sout := <-c.sChannel:
				teeOut.Write([]byte(sout))
			case eout := <-c.eChannel:
				teeErr.Write([]byte(eout))
			}
		}
	}(c)

	return c
}

// Stop does return to normal state
func Stop(c *Container) {

	// only if loggig was enabled
	if nil == c {
		return
	}

	// close pipes
	_ = c.writerStdout.Close()
	_ = c.writerStderr.Close()

	// wait
	time.Sleep(10 * time.Millisecond)

	// restore
	os.Stdout = c.origStdout
	os.Stderr = c.origStderr
	log.SetOutput(c.oldLog)

	// logfile
	c.lfHandle.Close()
}
