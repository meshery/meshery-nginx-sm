FROM golang:1.14-stretch as bd
ARG CONFIG_PROVIDER="viper"
RUN apt update && apt install git libc-dev gcc pkgconf -y
COPY ${PWD} /go/src/github.com/layer5io/meshery-nginx/
WORKDIR /go/src/github.com/layer5io/meshery-nginx/
RUN go build -ldflags="-w -s -X main.configProvider=$CONFIG_PROVIDER" -a -o meshery-nginx

FROM golang:1.14-stretch
RUN apt update && apt install ca-certificates curl -y
# Install kubectl
RUN curl -LO "https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl" && \
	chmod +x ./kubectl && \
	mv ./kubectl /usr/local/bin/kubectl

RUN mkdir /home/scripts/ && \
	mkdir -p ${HOME}/.kube/

COPY --from=bd /go/src/github.com/layer5io/meshery-nginx/meshery-nginx /home/
COPY ${PWD}/scripts /home/scripts
WORKDIR /home
CMD ./meshery-nginx
