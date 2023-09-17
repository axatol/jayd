FROM node:18-alpine as web
WORKDIR /app
COPY ./package.json ./package-lock.json ./
COPY ./web/package.json ./web/package.json
RUN npm install
COPY ./web ./web

ARG AUTH0_DOMAIN
ENV VITE_AUTH0_DOMAIN=${AUTH0_DOMAIN}
ARG AUTH0_CLIENT_ID
ENV VITE_AUTH0_CLIENT_ID=${AUTH0_CLIENT_ID}
ARG API_URL
ENV VITE_API_URL=${API_URL}
ARG API_AUDIENCE
ENV VITE_API_AUDIENCE=${API_AUDIENCE}
ARG COMMIT_SHA
ENV VITE_COMMIT_SHA=${COMMIT_SHA}
ARG NODE_ENV=production
ENV NODE_ENV=${NODE_ENV}

RUN npm run build

FROM golang:1.20-alpine as server
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG COMMIT_SHA=unknown
RUN go build \
  -o /bin/app \
  -ldflags="-X 'github.com/axatol/jayd/pkg/config.BuildCommit=${COMMIT_SHA}' -X 'github.com/axatol/jayd/pkg/config.BuildTime=$(date +"%Y-%m-%dT%H:%M:%S%z")'" \
  ./cmd/server/main.go

FROM python:3.11-slim-bullseye
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
  ca-certificates \
  curl \
  ffmpeg \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/download/2023.07.06/yt-dlp -o /usr/local/bin/yt-dlp \
  && chmod a+rx /usr/local/bin/yt-dlp

WORKDIR /go/app
RUN adduser --disabled-password --gecos "" --uid 1000 default
USER default
COPY --from=server --chown=default /bin/app /bin/app
COPY --from=web --chown=default /app/web/dist /web
ENTRYPOINT [ "/bin/app" ]
