package errs

import "net/http"

var write = http.ResponseWriter.Write

func fakeWrite(http.ResponseWriter, []byte) (int, error) {
	return 0, ErrResponseWriter
}

func restoreWrite(replace func(http.ResponseWriter, []byte) (int, error)) {
	write = replace
}
