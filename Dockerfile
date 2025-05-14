FROM alpine:3.21

# 1) create the user as root
RUN adduser -D appuser

# 2) copy down your statically-linked binary while you’re still root
COPY --from=builder /app/bin/main /usr/local/bin/

# 3) now switch to the unprivileged user
USER appuser

# 4) set a sensible working directory
WORKDIR /home/appuser

EXPOSE 8080

# 5) exec your binary
CMD ["/usr/local/bin/main"]
