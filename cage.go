package cage

import (
	"bytes"
	"io"
	"os"
	"time"
)

type Container struct {
	origStdout   *os.File // old
	writerStdout *os.File // new

	origStderr   *os.File // old
	writerStderr *os.File // new

	sChannel chan string // channel for os.Stdout messages
	eChannel chan string // channel for os.Stderr messages

	//data   string
	//Data []string

	lfHandle *os.File // logfile handle
}

// Start does append all output parallel into a logfile with name fn
func Start(fn string) *Container {

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

	// create container
	c := &Container{
		origStdout:   os.Stdout,
		writerStdout: wStdout,

		origStderr:   os.Stderr,
		writerStderr: wStderr,

		sChannel: make(chan string),
		eChannel: make(chan string),

		lfHandle: lfH,
	}

	// re-direct
	os.Stdout = c.writerStdout // all to os.Stdout goees now to c.writerStdout
	os.Stderr = c.writerStderr // all to os.Stderr goees now to c.writerStderr

	// writing to os.Stdout will go to pipe and come out of rStdout now
	// writing to os.Stderr will go to pipe and come out of rStderr now

	// create duplication
	teeOut := io.MultiWriter(c.origStdout, c.lfHandle)
	teeErr := io.MultiWriter(c.origStderr, c.lfHandle)

	// pipes reader
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
	_ = c.writerStdout.Close()
	_ = c.writerStderr.Close()
	time.Sleep(10 * time.Millisecond)

	os.Stdout = c.origStdout
	os.Stderr = c.origStderr

	c.lfHandle.Close()

	//c.Data = strings.Split(c.data, "\n")
	//if c.Data[len(c.Data)-1] == "" {
	//	c.Data = c.Data[:len(c.Data)-1]
	//}
}
