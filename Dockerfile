FROM alpine:3.6
WORKDIR /app
# Now just add the binary
COPY arMonitor /app/
ENTRYPOINT ["/app/arMonitor"]
EXPOSE 8002