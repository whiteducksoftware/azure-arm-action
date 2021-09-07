# Builder
FROM golang:1.16-alpine as builder
WORKDIR /app

# Add src files
ADD . .

# Fetch dependencies.
RUN go mod download
RUN go mod verify

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -o /usr/local/bin/azure-arm-action

# Runner
FROM gcr.io/distroless/base

# Copy our static executable.
COPY --from=builder /usr/local/bin/azure-arm-action /usr/local/bin/azure-arm-action

# Add lables
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="GitHub Action Deploy Azure ARM" 
LABEL org.label-schema.description="GitHub Action which can deploy Azure Resource Manager (ARM) templates" 
LABEL org.label-schema.vcs-ref="https://github.com/whiteducksoftware/azure-arm-action"
LABEL org.label-schema.maintainer="Stefan Kürzeder <stefan.kuerzeder@whiteduck.de>"

ENTRYPOINT ["/usr/local/bin/azure-arm-action"]