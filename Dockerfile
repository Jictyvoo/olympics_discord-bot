### STAGE 1: Build ###
FROM golang:1.22-bullseye AS builder

ENV PATH $GOPATH/bin:$PATH
ENV CGO_ENABLED=1
ENV GO1111MODULE=on


RUN apt update && \
    apt install -y --no-install-recommends \
    git \
    gcc \
    libc6-dev \
    pkg-config \
    libcurl4-openssl-dev && \
    rm -rf /var/lib/apt/lists/*


ENV SOURCE_DIR=/go/src/olympics_dicord-bot
WORKDIR $SOURCE_DIR

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
ADD go.* ./
RUN go mod download

# Copy all the Code and stuff to compile everything
ADD . .

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -a -installsuffix cgo -o olympicsBOT ./cmd/discordbot
#
#
########################################################################################################################
# Moving the binary to the 'final Image' to make it smaller
FROM debian:buster-slim

ENV DISCORD_TOKEN=''
ENV DISCORD_CLIENT_ID=''
ENV WATCH_COUNTRIES=''
ENV API_LOCALE='ENG'

RUN apt update && \
        apt install -y --no-install-recommends \
        ca-certificates curl \
        && rm -rf /var/lib/apt/lists/*

RUN curl -sSf https://atlasgo.sh | sh

WORKDIR /home/olympics_dicord-bot

# Copy the generated binary from builder image to execution image
COPY --from=builder /go/src/olympics_dicord-bot/olympicsBOT /bin/olympicsBOT

ADD build/migrations ./build/migrations
ADD build/container-entrypoint.sh ./entrypoint.sh

# Ensure the binary is executable
RUN chmod +x /bin/olympicsBOT && chmod +x ./entrypoint.sh


# Run the binary program
#ENTRYPOINT ["tail", "-f", "/dev/null"]
ENTRYPOINT ["/home/olympics_dicord-bot/entrypoint.sh"]
CMD ["/home/olympics_dicord-bot/entrypoint.sh"]
