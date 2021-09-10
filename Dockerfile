
FROM alpine:3.11.7 as intermediate

ARG SSH_PRIVATE_KEY
RUN apk add -qU openssh
RUN mkdir ~/.ssh/
RUN echo "${SSH_PRIVATE_KEY}" > ~/.ssh/id_ed25519
RUN chmod 600 ~/.ssh/id_ed25519
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts

FROM alpine:3.11.7

ENV GOPATH /go

RUN mkdir ~/.ssh/
COPY --from=intermediate /root/.ssh/id_ed25519 /root/.ssh/id_ed25519
COPY --from=intermediate /root/.ssh/known_hosts /root/.ssh/known_hosts

COPY . /go/src/github.com/hound-search/hound

RUN apk update \
	&& apk add go git subversion libc-dev mercurial bzr openssh tini \
	&& cd /go/src/github.com/hound-search/hound \
	&& go mod download \
	&& go install github.com/hound-search/hound/cmds/houndd \
	&& apk del go \
	&& rm -f /var/cache/apk/* \
	&& rm -rf /go/src /go/pkg

RUN mkdir ./data
COPY config.json ./data

EXPOSE 6080

ENTRYPOINT ["/sbin/tini", "--", "/go/bin/houndd", "-conf", "/data/config.json"]
