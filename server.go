package main

import (
	"archive/tar"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/codegangsta/negroni"
	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
)

var dock *docker.Client
var dockerAuth *docker.AuthConfigurations

func main() {
	var (
		dockerHost      = os.Getenv("DOCKER_HOST")
		dockerCertPath  = os.Getenv("DOCKER_CERT_PATH")
		dockerTlsVerify = os.Getenv("DOCKER_TLS_VERIFY") != ""
	)

	var (
		defaultCaFile   = "ca.pem"
		defaultKeyFile  = "key.pem"
		defaultCertFile = "cert.pem"
	)

	if dockerCertPath == "" {
		dockerCertPath = filepath.Join(os.Getenv("HOME"), ".docker")
	}

	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
	}

	var err error

	if dockerTlsVerify {
		dock, err = docker.NewTLSClient(dockerHost,
			filepath.Join(dockerCertPath, defaultCertFile),
			filepath.Join(dockerCertPath, defaultKeyFile),
			filepath.Join(dockerCertPath, defaultCaFile))
	} else {
		dock, err = docker.NewClient(dockerHost)
	}

	if err != nil || dock == nil {
		log.Fatal("couldn't initialise Docker", err)
	}

	if err := dock.Ping(); err != nil {
		log.Fatal("couldn't ping Docker", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/build", build).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":" + port)
}

func build(res http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()

	authConfigs, err := authsFromHeaders(req.Header)
	if err != nil {
		panic(err)
	}

	authConfig, err := authFromHeaders(req.Header)
	if err != nil {
		panic(err)
	}

	buildStream := toReader(func(w io.Writer) error {
		return addBuildpack(w, req.Body)
	})

	bodyFormatter := toWriter(func(r io.Reader) error {
		return formatJSON(res, r)
	})

	buildOpts := docker.BuildImageOptions{
		Name:          vars.Get("t"),
		InputStream:   buildStream,
		OutputStream:  bodyFormatter,
		RawJSONStream: true,
		NoCache:       true,
		AuthConfigs:   authConfigs,
	}

	if err := dock.BuildImage(buildOpts); err != nil {
		panic(err)
	}

	pushOpts := docker.PushImageOptions{
		Name:          buildOpts.Name,
		OutputStream:  bodyFormatter,
		RawJSONStream: true,
	}

	if err := dock.PushImage(pushOpts, authConfig); err != nil {
		panic(err)
	}
}

func addBuildpack(dest io.Writer, src io.ReadCloser) error {
	r := tar.NewReader(src)
	w := tar.NewWriter(dest)

	sawDockerfile := false
	filter := func(hdr *tar.Header) bool {
		if path.Clean(hdr.Name) == "Dockerfile" {
			sawDockerfile = true
		}
		return true
	}

	if err := copyTar(w, r, filter); err != nil {
		return err
	}

	if !sawDockerfile {
		if err := addFile(w, "packs/node/Dockerfile"); err != nil {
			return err
		}
	}

	return w.Flush()
}
