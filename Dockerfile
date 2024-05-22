FROM public.ecr.aws/docker/library/golang:latest AS build

WORKDIR /app
COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN GOARCH=amd64 GOOS=linux go build -v ./cmd/pagemail/

ARG BUCKET_NAME
ENV BUCKET_NAME=$BUCKET_NAME

CMD ["/app/pagemail"]
