# container-factory

  container-factory produces Docker images from tarballs of application source code.
  It accepts archives with `Dockerfile`s, but if your application's language is supported,
  it can automatically add a suitable `Dockerfile`.

  Currently, only node.js is supported, marked by the presence of a `package.json`.

## License

  container-factory is distributed under the terms of the ISC license.

## Installation

  Make sure you have [Docker](https://docker.io) set up.
  We access the Docker API, so we pass that through when running the container.
  The Docker server is autodetected in the same way as the Docker CLI does,
  and like the Docker CLI, we default to `/var/run/docker.sock`.

  Here's how to run container-factory in a container:
```
docker run -d -p 9001:3000 -v /var/run/docker.sock:/var/run/docker.sock lsqio/container-factory
```

## API
### POST /build

  Build a container. Takes an application tarball as body.
  Modeled on [Docker's `/build` API](https://docs.docker.com/reference/api/docker_remote_api_v1.18/#build-image-from-a-dockerfile)

  Parameters:
    
    * t: Docker tag to publish the resulting Docker image as.

### Response

  Docker-style progress reporting, with the push happening right after the build.

  Build progress:
```json
{"stream": "Step 1..."}
{"stream": "..."}
{"error": "Error...", "errorDetail": {"code": 123, "message": "Error..."}}
```

  Push progress:
```json
{"status": "Pushing..."}
{"status": "Pushing", "progress": "1/? (n/a)", "progressDetail": {"current": 1}}}
{"error": "Invalid..."}
```

