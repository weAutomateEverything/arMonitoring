FROM alpine:3.6

WORKDIR /app
# Now just add the binary
COPY arMonitoring /app/
ENTRYPOINT ["/app/arMonitoring"]
EXPOSE 8002