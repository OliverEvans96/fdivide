FROM golang AS build
COPY . /go/src/github.com/OliverEvans96/fdivide
WORKDIR /go/src/github.com/OliverEvans96/fdivide
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install ./...

FROM alpine AS run
ARG GIT_COMMIT
ENV GIT_COMMIT=$GIT_COMMIT
COPY --from=build /go/bin/ /usr/bin
