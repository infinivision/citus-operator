#!/bin/bash

#https://www.postgresql.org/docs/current/libpq-pgpass.html
# prepare pgpass:
# echo "*:*:*:postgres:supassword" > pgpass && chmod 600 pgpass && sudo chown haproxy pgpass && sudo mv pgpass /etc/haproxy/pgpass

PGPASSFILE=/etc/haproxy/pgpass psql --host "$HAPROXY_SERVER_ADDR" --port "$HAPROXY_SERVER_PORT" --username=postgres --no-password --pset=tuples_only -c "SELECT pg_is_in_recovery();" | grep -P "^\s*f\s*$" > /dev/null
RET=$?
echo `date --rfc-3339=seconds` "$HAPROXY_SERVER_ADDR:$HAPROXY_SERVER_PORT pg_is_master RET $RET"
exit "$RET"
