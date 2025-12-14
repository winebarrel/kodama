FROM golang:1.25 AS build

WORKDIR /src
COPY go.* /src/
RUN go mod download

COPY *.go /src/
COPY cmd/kodama/*.go /src/cmd/kodama/
ARG KODAMA_VERSION
RUN CGO_ENABLED=0 go build -o kodama -ldflags "-X main.version=${KODAMA_VERSION}" ./cmd/kodama

FROM gcr.io/distroless/static

COPY --from=build /src/kodama /

ENTRYPOINT ["/kodama"]
