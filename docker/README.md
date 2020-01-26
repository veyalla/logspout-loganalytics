# Docker
## How to build this image

### Configure the multiarch builder with experimental buildx.
[https://docs.docker.com/engine/reference/commandline/cli/#experimental-features](https://docs.docker.com/engine/reference/commandline/cli/#experimental-features)

### Create a buildx environment
```sh
docker buildx create --name multiarch-builder
docker buildx use multiarch-builder
docker buildx inspect --bootstrap
```

### Launch the build
```sh
docker buildx build --push --platform linux/arm/v7,linux/amd64,linux/arm64,linux/386 -t kbeaugrand/logspout-loganalytics -t kbeaugrand/logspout-loganalytics:0.0.1 .
```