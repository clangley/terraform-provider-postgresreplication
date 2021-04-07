FROM postgres:12-alpine
COPY postgresql.conf /etc/postgresql/postgresql.conf
RUN apk add --no-cache git build-base clang llvm
RUN git clone https://github.com/eulerto/wal2json.git && \
    cd wal2json && \
    USE_PGXS=1 make && \
    USE_PGXS=1 make install

COPY entrypoint.sql /docker-entrypoint-initdb.d/entrypoint.sql

CMD ["-c", "config_file=/etc/postgresql/postgresql.conf"]
