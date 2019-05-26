FROM postgres:11 AS pg_builder

# install Citus
RUN apt-get update \
    && apt-get install -y postgresql-server-dev-$PG_MAJOR wget build-essential

RUN wget https://github.com/ChenHuajun/pg_roaringbitmap/archive/v0.2.1.tar.gz \
    && tar xzvf v0.2.1.tar.gz \
    && cd pg_roaringbitmap-0.2.1 \
    && make \
    && make install \
    && tar czvf /postgresql-$PG_MAJOR-roaringbitmap-0.2.1.tgz /usr/lib/postgresql/$PG_MAJOR/lib/bitcode/roaringbitmap* /usr/lib/postgresql/$PG_MAJOR/lib/roaringbitmap.so /usr/share/postgresql/$PG_MAJOR/extension/roaringbitmap*

RUN apt-get install -y protobuf-c-compiler libprotobuf-c0-dev \
    && wget https://github.com/citusdata/cstore_fdw/archive/v1.6.2.tar.gz \
    && tar xzvf v1.6.2.tar.gz \
    && cd cstore_fdw-1.6.2 \
    && make \
    && make install \
    && tar czvf /postgresql-$PG_MAJOR-cstore_fdw-1.6.2.tgz /usr/lib/postgresql/$PG_MAJOR/lib/bitcode/cstore_fdw* /usr/lib/postgresql/$PG_MAJOR/lib/cstore_fdw.so /usr/share/postgresql/$PG_MAJOR/extension/cstore_fdw*

RUN wget https://github.com/sorintlab/stolon/releases/download/v0.13.0/stolon-v0.13.0-linux-amd64.tar.gz \
    && tar -C /root/ -xzvf stolon-v0.13.0-linux-amd64.tar.gz

RUN wget -O /root/gosu https://github.com/tianon/gosu/releases/download/1.11/gosu-amd64 \
    && chmod +x /root/gosu

FROM postgres:11 AS keeper
ARG VERSION=8.2.1
LABEL maintainer="Citus Data https://citusdata.com" \
    org.label-schema.name="Citus" \
    org.label-schema.description="Scalable PostgreSQL for multi-tenant and real-time workloads" \
    org.label-schema.url="https://www.citusdata.com" \
    org.label-schema.vcs-url="https://github.com/citusdata/citus" \
    org.label-schema.vendor="Citus Data, Inc." \
    org.label-schema.version=${VERSION} \
    org.label-schema.schema-version="1.0"

ENV CITUS_VERSION ${VERSION}.citus-1

# install Citus
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    && curl -s https://install.citusdata.com/community/deb.sh | bash \
    && apt-get install -y postgresql-$PG_MAJOR-citus-8.2=$CITUS_VERSION \
    postgresql-$PG_MAJOR-hll=2.12.citus-1 \
    postgresql-$PG_MAJOR-topn=2.2.0 \
    postgresql-$PG_MAJOR-cron=1.1.4-1.pgdg90+1 \
    && apt-get purge -y --auto-remove curl \
    && rm -rf /var/lib/apt/lists/*

# install roaringbitmap
COPY --from=pg_builder /usr/lib/postgresql/$PG_MAJOR/lib/bitcode/roaringbitmap* /usr/lib/postgresql/$PG_MAJOR/lib/bitcode/
COPY --from=pg_builder /usr/lib/postgresql/$PG_MAJOR/lib/roaringbitmap.so /usr/lib/postgresql/$PG_MAJOR/lib/
COPY --from=pg_builder /usr/share/postgresql/$PG_MAJOR/extension/roaringbitmap* /usr/share/postgresql/$PG_MAJOR/extension/

# install cstore_fdw
COPY --from=pg_builder /usr/lib/postgresql/$PG_MAJOR/lib/bitcode/cstore_fdw* /usr/lib/postgresql/$PG_MAJOR/lib/bitcode/
COPY --from=pg_builder /usr/lib/postgresql/$PG_MAJOR/lib/cstore_fdw.so /usr/lib/postgresql/$PG_MAJOR/lib/
COPY --from=pg_builder /usr/share/postgresql/$PG_MAJOR/extension/cstore_fdw* /usr/share/postgresql/$PG_MAJOR/extension/
COPY --from=pg_builder /usr/lib/x86_64-linux-gnu/libprotobuf-c.so.1 /usr/lib/x86_64-linux-gnu/libprotobuf-c.so.1

# add citus to default PostgreSQL config
RUN echo "shared_preload_libraries='citus, pg_cron, cstore_fdw'" >> /usr/share/postgresql/postgresql.conf.sample

# install stolon-keeper
COPY --from=pg_builder /root/stolon-v0.13.0-linux-amd64/bin/stolon-keeper /usr/local/bin/

# install gosu
COPY --from=pg_builder /root/gosu /usr/local/bin/

ENTRYPOINT chown -R postgres:postgres /stolon-data && exec gosu postgres:postgres /usr/local/bin/stolon-keeper --data-dir /stolon-data --pg-bin-path /usr/lib/postgresql/11/bin --pg-listen-address 0.0.0.0


FROM ubuntu AS proxy

# install stolon-proxy
COPY --from=pg_builder /root/stolon-v0.13.0-linux-amd64/bin/stolon-proxy /usr/local/bin/

ENTRYPOINT /usr/local/bin/stolon-proxy


FROM ubuntu AS sentinel

# install stolon-sentinel
COPY --from=pg_builder /root/stolon-v0.13.0-linux-amd64/bin/stolon-sentinel /usr/local/bin/

ENTRYPOINT /usr/local/bin/stolon-sentinel
