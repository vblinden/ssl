FROM golang:1.22-alpine as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o bin ./

FROM gcr.io/distroless/base-debian11
WORKDIR /

COPY --from=builder /usr/src/app/bin /bin

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["bin"]
