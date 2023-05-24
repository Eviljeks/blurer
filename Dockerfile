FROM golang:1.19-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY cmd cmd
COPY internal internal
COPY migrations migrations
COPY pkg pkg

RUN go mod download

RUN go build -o /blurer ./cmd/server/

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /blurer /blurer

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/blurer"]