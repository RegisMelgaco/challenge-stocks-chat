FROM golang:1.18 as build

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api  cmd/api/main.go


FROM ubuntu:22.04

RUN apt-get update -y

WORKDIR /app/

COPY --from=build /go/src/app/api .

RUN useradd --create-home --uid 1000 gopher
USER 1000

CMD ["./api"]
