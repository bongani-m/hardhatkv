FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/hardhatkv .

FROM debian:bookworm-slim
RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates \
	&& rm -rf /var/lib/apt/lists/*
COPY --from=build /out/hardhatkv /usr/local/bin/hardhatkv
COPY docker/redis.conf /etc/hardhatkv/redis.conf
ENV CONFIG=/etc/hardhatkv/redis.conf
WORKDIR /data
VOLUME /data
EXPOSE 6399
ENTRYPOINT ["hardhatkv"]
