#!/bin/sh
# wait-for.sh

echo "start waiting for api..."

if [ $# -ne 3 ]; then
   echo "provide 3 arguments: api url, timeout (s) and callback script"
   exit 0
fi

for i in $(seq 1 $2)
do
  out=$(curl -sS "$1/health") || echo "waiting for api connection ..."

  if [ "$out" = "{\"status\":\"up\"}" ]
  then
    echo "connected to account api"
    sh "$3" $1
    exit 1
  fi

  sleep 1
done


echo "timeout"
exit 0


#
#>&2 echo "Postgres is up - executing command"
#exec "$@"