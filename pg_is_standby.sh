#!/bin/bash

#https://www.postgresql.org/docs/current/libpq-pgpass.html
# prepare pgpass:
# echo "*:*:*:postgres:supassword" > pgpass && chmod 600 pgpass && sudo chown haproxy pgpass && sudo mv pgpass /etc/haproxy/pgpass

# psql checking is heavy, each run costs 300ms~1s.
#PGPASSFILE=/etc/haproxy/pgpass psql "host=$HAPROXY_SERVER_ADDR port=$HAPROXY_SERVER_PORT connect_timeout=1" --username=postgres --no-password --pset=tuples_only -c "SELECT pg_is_in_recovery();" | grep -P "^\s*t\s*$" > /dev/null

# Assumes keeper metric port is postgres port plus 1.
let STKEEPER_METRICS_PORT=$HAPROXY_SERVER_PORT+1
curl http://$HAPROXY_SERVER_ADDR:$STKEEPER_METRICS_PORT/metrics 2>/dev/null | grep -P "^stolon_keeper_local_role.*standby.*1$" > /dev/null

RET=$?
#echo `date --rfc-3339=seconds` "$HAPROXY_SERVER_ADDR:$HAPROXY_SERVER_PORT pg_is_standby RET $RET"
exit "$RET"
