package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestTrip(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		input    []flight
		expected flight
	}{
		{[]flight{{"SFO", "EWR"}}, flight{"SFO", "EWR"}},
		{[]flight{{"SFO", "ATL"}, {"ATL", "EWR"}}, flight{"SFO", "EWR"}},
		{[]flight{{"ATL", "EWR"}, {"SFO", "ATL"}}, flight{"SFO", "EWR"}},
		{[]flight{{"SFO", "ATL"}, {"ATL", "GSO"}, {"GSO", "EWR"}}, flight{"SFO", "EWR"}},
		{[]flight{{"GSO", "EWR"}, {"ATL", "GSO"}, {"SFO", "ATL"}}, flight{"SFO", "EWR"}},
		{[]flight{{"GSO", "EWR"}, {"SFO", "ATL"}, {"ATL", "GSO"}}, flight{"SFO", "EWR"}},
		{[]flight{{"IND", "EWR"}, {"SFO", "ATL"}, {"GSO", "IND"}, {"ATL", "GSO"}}, flight{"SFO", "EWR"}},
	} {
		if a, e := trip(test.input), test.expected; a != e {
			t.Errorf("%#v: actual %#v, expected %#v", test.input, a, e)
		}
	}
}

func TestHandler(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest("POST", "/calculate", bytes.NewBufferString(`{"data":{"flights":[["ATL","EWR"],["SFO","ATL"]]}}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	rec := httptest.NewRecorder()
	handler().ServeHTTP(rec, req)
	if a, e := rec.Code, http.StatusOK; a != e {
		t.Errorf("expected %#v, actual %#v", a, e)
	}
	var res response
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}
	if res.Data == nil {
		t.Fatal("data is nil")
	}
	tr, ok := res.Data["trip"]
	if !ok {
		t.Fatal("trip is nil")
	}
	if a, e := tr, ([]any{"SFO", "EWR"}); !reflect.DeepEqual(a, e) {
		t.Errorf("expected %#v, actual %#v", a, e)
	}
}
