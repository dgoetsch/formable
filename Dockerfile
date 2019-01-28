FROM hashicorp/terraform:0.11.7 as terraform

FROM golang:1.10 as gobuild
WORKDIR /go/src/github.com/dgoetsch/formable
COPY . .
RUN go build -o formable .

FROM google/cloud-sdk:214.0.0-alpine
COPY --from=terraform /bin/terraform /bin/
COPY --from=gobuild /go/src/github.com/dgoetsch/formable/formable /bin/
COPY entrypoint.sh /app/

ENV GOOGLE_APPLICATION_CREDENTIALS=/mnt/identity.json
ENV CONFIG_DIR=/mnt/config
ENV TF_DIR=/mnt/terraform
ENV TF_CMD=plan

ENTRYPOINT ["/bin/sh", "/app/entrypoint.sh"]
