package cage

import (
	"io"
	"log"
	"os"
)

// Container keeps re-direction informantion
type Container struct {
	// old
	oldLog     io.Writer
	origStdout *os.File
	origStderr *os.File

	// new
	writerStdout *os.File
	writerStderr *os.File
	lfHandle     *os.File // logfile handle
	lfName       string
}

// Start does append all output parallel into a logfile with name fn
func Start(fn string) *Container {

	// start logging only if fn not "none"
	if "none" == fn {
		log.Println("No logfile writing...")
		return nil
	}

	// open logfile
	lfH, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file %s: %v", fn, err)
		fn = "off"
		return nil
	}
	log.Printf("Writing to logfile %s...\n", fn)

	// open pipes
	rStdout, wStdout, _ := os.Pipe()
	rStderr, wStderr, _ := os.Pipe()

	// create container for recovering
	c := &Container{
		oldLog: log.Writer(),

		origStdout: os.Stdout,
		origStderr: os.Stderr,

		writerStdout: wStdout,
		writerStderr: wStderr,

		lfHandle: lfH,
		lfName:   fn,
	}

	// re-direct
	log.SetOutput(io.MultiWriter(c.oldLog, c.lfHandle)) // writing to log will go also to logfile now
	os.Stdout = c.writerStdout                          // all to os.Stdout goees now to c.writerStdout and comee out of rStdout now
	os.Stderr = c.writerStderr                          // all to os.Stderr goees now to c.writerStderr and comes out of rStderr now

	// create duplication
	teeOut := io.MultiWriter(c.origStdout, c.lfHandle)
	teeErr := io.MultiWriter(c.origStderr, c.lfHandle)

	// copy from pipe to tee
	go func(w io.Writer, r io.Reader) {
		io.Copy(w, r)
	}(teeOut, rStdout)
	go func(w io.Writer, r io.Reader) {
		io.Copy(w, r)
	}(teeErr, rStderr)

	return c
}

// Stop does return to normal state
func Stop(c *Container) {

	// only if loggig was enabled
	if nil == c {
		log.Println("No logfile writing...done")
		return
	}

	// close pipes
	_ = c.writerStdout.Close()
	_ = c.writerStderr.Close()

	// restore
	os.Stdout = c.origStdout
	os.Stderr = c.origStderr
	log.SetOutput(c.oldLog)

	// logfile
	c.lfHandle.Close()
	log.Printf("Writing to logfile %s...done\n", c.lfName)
}
