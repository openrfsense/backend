FROM scratch

COPY config.yml /config.yml
COPY orfs-backend /orfs-backend

ENTRYPOINT ["/orfs-backend"]