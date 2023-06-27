FROM golang:1.20.5 AS builder

ENV GO111MODULE=on
WORKDIR /app

RUN mkdir /var/log/entry/

# The <src> path must be inside the context of the build
# you cannot COPY ../something /something 
# because the first step of a docker build is to send the 
# context directory (and subdirectories) to the docker daemon.
#
# In below command, the context directory is ~/projects/social-network
# In specific, we will call: `docker build -t app_svc.linux:latest -f cmd/authen_and_post_svc/Dockerfile .`
# from ~/projects/social-network
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/cache go mod download
RUN --mount=type=cache,target=/go/pkg/mod/cache CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o aap_svc.linux -mod=mod cmd/authen_and_post_svc/main.go
CMD ["/app/aap_svc.linux"]

FROM builder AS test_env
RUN cp /app/configs/files/test.yml /app/config.yml
CMD ["/app/aap_svc.linux"]

FROM builder AS live_env
RUN cp /app/configs/files/live.yml /app/config.yml
CMD ["/app/aap_svc.linux"]