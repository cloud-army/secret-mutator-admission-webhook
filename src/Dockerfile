FROM alpine:latest
RUN apk --no-cache add ca-certificates && \
  update-ca-certificates

COPY webhook /bin/webhook
ENTRYPOINT ["/bin/webhook"]
