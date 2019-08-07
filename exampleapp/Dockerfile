FROM ubuntu 
MAINTAINER Anubhav Mishra <mishra@hashicorp.com>

ENV VAULT_VERSION=1.1.3

RUN apt-get update
RUN apt-get install -y gnupg2 curl git unzip bash libpcap-dev openssl

RUN curl -LOks https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip && \
    unzip vault_${VAULT_VERSION}_linux_amd64.zip -d /tmp/vault_${VAULT_VERSION} && \
    mv /tmp/vault_${VAULT_VERSION}/vault /usr/local/bin/vault && \
    rm vault_${VAULT_VERSION}_linux_amd64.zip

ADD build/linux/amd64/exampleapp /usr/bin/exampleapp
ADD token.sh /usr/bin/token.sh
ADD wrapper.sh /usr/bin/wrapper.sh
ENTRYPOINT ["wrapper.sh"]
