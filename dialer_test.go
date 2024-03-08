package geario

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"net"
	"time"
	"io"
)

func TestDialer(t *testing.T) {
	svc := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("OK"))
	}))

	cli := svc.Client()

	cli.Transport = &http.Transport{
		DialContext: GearDialer(&net.Dialer{}, time.Second, 50).DialContext,
	}

	req, err := http.NewRequest("GET", "http://"+svc.Listener.Addr().String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	resp, err := cli.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "OK" {
		t.Fatal("want OK, got", string(data))
	}

	if time.Since(start) < 3*time.Second {
		t.Fatal("time too short")
	}
}
