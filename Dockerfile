FROM scratch
COPY config.yml /
ENTRYPOINT ["/orfs-backend"]