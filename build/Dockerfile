# TODO: Auto-update version?
FROM registry.gitlab.com/lavo-nutrition/graphql-schema/yoga-server:3dc333ef AS prisma

FROM registry.gitlab.com/lavo-nutrition/graphql-translator/go-get:latest AS build
COPY --from=prisma /tmp/mypkg/prisma/generated/go-client /go/src/gitlab.com/lavo-nutrition/graphql-schema/prisma
COPY . /go/src/gitlab.com/lavo-nutrition/graphql-translator
WORKDIR /go/src/gitlab.com/lavo-nutrition/graphql-translator
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install ./...

FROM alpine AS run
COPY --from=build /go/bin/ /usr/bin
