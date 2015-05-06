package main

import (
	"archive/tar"
	"io"
	"os"
)

func addFile(w *tar.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}

	if err := w.WriteHeader(hdr); err != nil {
		return err
	}

	if fi.Mode().IsRegular() {
		if _, err := io.Copy(w, f); err != nil {
			return err
		}
	}

	return nil
}

func copyTar(dest *tar.Writer, src *tar.Reader, f func(*tar.Header) bool) error {
	for {
		hdr, err := src.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if f != nil && !f(hdr) {
			continue
		}

		if err := dest.WriteHeader(hdr); err != nil {
			return err
		}

		if _, err := io.Copy(dest, src); err != nil {
			return err
		}
	}

	return nil
}
