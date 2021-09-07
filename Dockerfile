# Builder
FROM golang:1.16-alpine as builder
WORKDIR /app

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && \
    apk add --no-cache git ca-certificates && \
    update-ca-certificates

# Add src files
ADD . .

# Fetch dependencies.
RUN go mod download
RUN go mod verify

# Build the binary.
ARG GIT_SHA
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -o /usr/local/bin/azure-arm-action

# Runner
FROM scratch

# Import the user and group files from the builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable.
COPY --from=builder /usr/local/bin/azure-arm-action /usr/local/bin/azure-arm-action

# Add lables
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="GitHub Action Deploy Azure ARM" 
LABEL org.label-schema.description="GitHub Action which can deploy Azure Resource Manager (ARM) templates" 
LABEL org.label-schema.vcs-ref="https://github.com/whiteducksoftware/azure-arm-action"
LABEL org.label-schema.maintainer="Stefan Kürzeder <stefan.kuerzeder@whiteduck.de>"

ENTRYPOINT ["/usr/local/bin/azure-arm-action"]