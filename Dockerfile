# Production Docker image, built as multistage build
# To build it, run:
#
#    docker build -t go-fractgen .
#
# To create a multi-platform build, use buildx:
# 
#    docker buildx build --platform linux/arm64,linux/amd64 -t registry.alexi.ch/go-fractgen:latest --push .
#
# Note that you need to have a compatible BuildKit builder. If you're using OrbStack, you can create one with:
#
# Create a parallel multi-platform builder
#    docker buildx create --name multiplatform-builder --use
# Make "buildx" the default
#    docker buildx install

FROM debian:bookworm-slim AS build

RUN apt-get update && apt-get install -y --no-install-recommends \
	golang-go \
	binutils \
	ca-certificates \
	curl \
	python3

RUN	go install golang.org/dl/go1.24.0@latest && \
	/root/go/bin/go1.24.0 download

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.2/install.sh | bash
RUN bash -c 'source /root/.nvm/nvm.sh && nvm install 22'

WORKDIR /src
RUN mkdir -p /build/webroot
COPY . .

# build go binary
RUN /root/go/bin/go1.24.0 build -o /build/go-fractgen main.go

# build web app
RUN bash -c 'source /root/.nvm/nvm.sh && \
	cd webroot && \
	nvm use && \
	npm install && \
	npm run build && \
	cp -r dist/* /build/webroot/'

############
FROM alpine:latest AS prod

EXPOSE 8000
WORKDIR /app
RUN adduser --uid 5000 --disabled-password fractgen
COPY --from=build /build /app
USER fractgen
ENTRYPOINT ["/app/go-fractgen"]
CMD [ "serve", "--listen=:8000" ]
