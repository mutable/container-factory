package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/fsouza/go-dockerclient"
)

func toReader(f func(io.Writer) error) io.Reader {
	r, w := io.Pipe()
	go func() {
		w.CloseWithError(f(w))
	}()

	return r
}

func toWriter(f func(io.Reader) error) io.Writer {
	r, w := io.Pipe()

	go func() {
		r.CloseWithError(f(r))
	}()

	return w
}

func formatJSON(dst io.Writer, src io.Reader) (err error) {
	decoder := json.NewDecoder(src)
	encoder := json.NewEncoder(dst)

	flusher, _ := dst.(http.Flusher)

	for {
		var data interface{}

		if err = decoder.Decode(&data); err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

		if err = encoder.Encode(&data); err != nil {
			return
		}

		if flusher != nil {
			flusher.Flush()
		}
	}
}

func authFromHeaders(headers map[string][]string) (auth docker.AuthConfiguration, err error) {
	for _, header := range headers["X-Registry-Auth"] {
		data, err := base64.URLEncoding.DecodeString(header)
		if err != nil {
			return auth, err
		}

		if err := json.Unmarshal(data, &auth); err != nil {
			return auth, err
		}

		return auth, nil
	}

	return
}

func authsFromHeaders(headers map[string][]string) (auth docker.AuthConfigurations, err error) {
	for _, header := range headers["X-Registry-Config"] {
		data, err := base64.URLEncoding.DecodeString(header)
		if err != nil {
			return auth, err
		}

		if err := json.Unmarshal(data, &auth); err != nil {
			return auth, err
		}

		return auth, nil
	}

	return
}
