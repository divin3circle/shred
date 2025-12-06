# GoReleaser builds the binary, so we just need a minimal runtime image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the pre-built binary (GoReleaser will provide this)
COPY shred .

# Make binary executable
RUN chmod +x shred

# Set entrypoint
ENTRYPOINT ["./shred"]

