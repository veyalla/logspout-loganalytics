# Docker
## How to build this image

```sh
docker buildx build -t kbeaugrand/logspout-loganalytics -t kbeaugrand/logspout-loganalytics:0.0.1 . -f Dockerfile --platform linux/arm/v7,linux/arm/v6,linux/amd64,linux/arm64,linux/ppc64le,linux/s390x,linux/386
```