FROM golang:1.19 as build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o /bin/app ./cmd/server/main.go

FROM python:3.11-slim-bullseye
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
  ca-certificates \
  curl \
  ffmpeg \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

RUN curl -L https://github.com/ytdl-patched/yt-dlp/releases/download/2023.02.17.334/yt-dlp -o /usr/local/bin/yt-dlp \
  && chmod a+rx /usr/local/bin/yt-dlp

WORKDIR /go/app
RUN adduser --disabled-password --gecos "" --uid 1000 default
USER default
COPY --from=build --chown=default /bin/app /bin/app
ENTRYPOINT [ "/bin/app" ]
