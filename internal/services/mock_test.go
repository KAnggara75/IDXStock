package services

import "net/http"

type roundTripperFunc struct {
	roundTrip func(req *http.Request) (*http.Response, error)
}

func (f *roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f.roundTrip(req)
}
