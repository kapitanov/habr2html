FROM golang:1.18-alpine AS backend

WORKDIR /src
COPY go.mod /src/go.mod
COPY go.sum /src/go.sum
ENV CGO_ENABLED 0
RUN go mod download

COPY . /src
RUN mkdir -p /out
RUN go test ./...
RUN go build -o /out/habr2html -o /out/habr2html -buildvcs=false ./cmd/habr2html

FROM alpine:3
WORKDIR /opt/habr2html
COPY --from=backend /out/habr2html /opt/habr2html/habr2html
ENTRYPOINT [ "/opt/habr2html/habr2html" ]
