package cage

import (
	"fmt"
	"os"
	"testing"
)

func TestStart(t *testing.T) {
	c := Start("logfile.txt")

	fmt.Println("test")
	fmt.Println("test2")
	fmt.Fprintln(os.Stderr, "stderr error")

	Stop(c)
	/*

		test
		test2
		stderr error


				fmt.Println(c.Data)

				if len(c.Data) != 3 {
					t.Error("Data length should be 3")
				}
				if c.Data[0] != "test" {
					t.Errorf("First line should be 'test', instead of %s", c.Data[0])
				}
				if c.Data[1] != "test2" {
					t.Errorf("Second line should be 'test2', instead of %s", c.Data[1])
				}
				if c.Data[2] != "stderr error" {
					t.Errorf("Third line should be 'stderr error', instead of %s", c.Data[2])
				}*/
}
