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
	b, err := FromBytesSize("120.912123123412311234123KiB")
	t.Log(b, err)

	a := NewBPSAver(time.Second)
	a.Add(b)
	t.Log(a.String())
}

func TestGearWriter(t *testing.T) {
	out := GearWriter(os.Stderr, time.Second, 5)
	begin := time.Now()
	out.Write([]byte("hello world\n"))
	end := time.Now()
	if end.Sub(begin) < time.Second*2 {
		t.Fail()
	}
}

func TestGearReader(t *testing.T) {
	out := GearReader(bytes.NewBufferString("hello world\n"), time.Second, 5)
	begin := time.Now()
	ioutil.ReadAll(out)
	end := time.Now()
	if end.Sub(begin) < time.Second*2 {
		t.Fail()
	}
}

func TestGearReadWriter(t *testing.T) {
	out := GearReadWriter(bytes.NewBuffer(nil), time.Second, 10)
	{
		begin := time.Now()
		fmt.Fprintf(out, "hello world\n")
		end := time.Now()
		if end.Sub(begin) < time.Second*1 {
			t.Fatalf("write time: %s", end.Sub(begin))
		}
	}

	{
		begin := time.Now()
		ioutil.ReadAll(out)
		end := time.Now()
		if end.Sub(begin) < time.Second*1 {
			t.Fatalf("read time: %s", end.Sub(begin))
		}
	}
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
