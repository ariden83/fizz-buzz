#Meant for running lint
ARG GOBUILDER_IMAGE
FROM ${GOBUILDER_IMAGE} as linter
ARG PROJECT_ROOT
WORKDIR $PROJECT_ROOT
COPY vendor vendor
COPY . .
RUN gometalinter.v1 ./... --fast --vendor --disable gotype --exclude=^benches/model/ --errors

#Meant for running swagger
ARG GOBUILDER_IMAGE
FROM golang:1.13.0-alpine3.10 as swagger

RUN apk add build-base openssh git

ARG PROJECT_ROOT
WORKDIR $PROJECT_ROOT
COPY vendor vendor
COPY . .

RUN export GO111MODULE=on

RUN go install -mod vendor github.com/go-swagger/go-swagger/cmd/swagger

RUN swagger generate spec -o ./swagger/swagger-template.json ./cmd ./...

#Meant for running tests
ARG GOBUILDER_IMAGE
FROM ${GOBUILDER_IMAGE} as tests
ARG PROJECT_ROOT
WORKDIR $PROJECT_ROOT
COPY --from=linter $PROJECT_ROOT/vendor ./vendor
COPY --from=swagger $PROJECT_ROOT/swagger ./swagger
COPY . .
ENV TZ="Europe/Paris"
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime
RUN echo $TZ > /etc/timezone

RUN go test ./... -v
RUN go test -run=XXX -bench=. ./...  -v
#RUN cat test.log

#Meant to build binary
ARG GOBUILDER_IMAGE
FROM ${GOBUILDER_IMAGE} as builder
ARG PROJECT_ROOT
ARG GIT_TAG_NAME
WORKDIR $PROJECT_ROOT
COPY --from=linter $PROJECT_ROOT/vendor ./vendor
COPY --from=linter $PROJECT_ROOT/pkg/api ./pkg/api
COPY --from=swagger $PROJECT_ROOT/swagger/ ./swagger/
COPY . .
RUN go build -ldflags "-X main.Version=$GIT_TAG_NAME" -o bin/main ./.

#Meant for building the deployment container
FROM alpine:3.10.1
ARG PROJECT_ROOT
WORKDIR /root
COPY --from=builder $PROJECT_ROOT/bin ./
COPY --from=swagger $PROJECT_ROOT/swagger/ ./swagger/
ENTRYPOINT ["./main"]
