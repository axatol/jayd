FROM node:18 as web
WORKDIR /app
COPY ./package.json ./package-lock.json ./
COPY ./web/package.json ./web/package.json
RUN npm install
COPY ./web ./web

ARG AUTH0_DOMAIN
ARG AUTH0_CLIENT_ID
ARG API_URL
ARG API_AUDIENCE
ENV VITE_AUTH0_DOMAIN=${AUTH0_DOMAIN}
ENV VITE_AUTH0_CLIENT_ID=${AUTH0_CLIENT_ID}
ENV VITE_API_URL=${API_URL}
ENV VITE_API_AUDIENCE=${API_AUDIENCE}

RUN npm run build

FROM golang:1.19 as server
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
COPY --from=server --chown=default /bin/app /bin/app
COPY --from=web --chown=default /app/web/dist /web
ENTRYPOINT [ "/bin/app" ]
