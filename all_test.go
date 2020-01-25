package geario

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	b, err := Parse("120.912123123412311234123KiB")
	t.Log(b, err)

	a := NewBPSAver(time.Second)
	a.Add(b)
	t.Log(a.String())
}

func TestGearWriter(t *testing.T) {
	out := GearWriter(os.Stdout, time.Second, 2)
	fmt.Fprintf(out, "hello world")
	fmt.Fprintf(out, "hello world")
}

func TestGearReader(t *testing.T) {
	out := GearReader(bytes.NewBufferString("hello world\nhello world"), time.Second, 2)
	ioutil.ReadAll(out)
}

func TestGearReadWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	out := GearReadWriter(buf, time.Second, 2)
	fmt.Fprintf(out, "hello world")
	fmt.Fprintf(out, "hello world")
	ioutil.ReadAll(out)
}

var bb = NewBPSAver(time.Second)

func TestGearSum(t *testing.T) {

	for i := 0; i != 10; i++ {
		time.Sleep(time.Second / 10)
		bb.Add(KiB)
	}
	t.Log(bb.String())
	for i := 0; i != 10; i++ {
		time.Sleep(time.Second / 10)
		bb.Add(KiB)
	}
	t.Log(bb.String())
}
