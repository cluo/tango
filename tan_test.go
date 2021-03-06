package tango

import (
	"reflect"
	"testing"
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"time"
	"io/ioutil"
)

func TestTan1(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	l := NewLogger(os.Stdout)
	o := Classic(l)
	o.Get("/", func() string {
		return Version()
	})
	o.Logger().Debug("it's ok")

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), Version())
}

func TestTan2(t *testing.T) {
	o := Classic()
	o.Get("/", func() string {
		return Version()
	})
	go o.Run()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8000/")
	if err != nil {
		t.Error(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	expect(t, resp.StatusCode, http.StatusOK)
	expect(t, string(bs), Version())
}

func TestTan3(t *testing.T) {
	o := Classic()
	o.Get("/", func() string {
		return Version()
	})
	go o.Run(":4040")

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:4040/")
	if err != nil {
		t.Error(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	expect(t, resp.StatusCode, http.StatusOK)
	expect(t, string(bs), Version())
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}