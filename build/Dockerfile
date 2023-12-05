FROM busybox:1.35

LABEL maintainer="R&D Fortuna a.s."

ADD syncer /

# Forward request and error logs to docker log collector
RUN chmod +x /syncer

EXPOSE 8080

ENTRYPOINT ["/syncer", "-http-port", "8080"]
