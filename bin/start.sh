#!/bin/sh

set -euo pipefail

if [ -z "${APP_NAME:-}" ] ; then
  echo "must set APP_NAME to start your service (example: APP_NAME=tigerhall-kittens-server bin/start.sh)"
  exit 1
fi

if [ ! -z "${POSTGRES_HOST:-}" ] && [ ! -z "${POSTGRES_PORT:-}" ]; then
  echo "wait for postgres @" $POSTGRES_HOST:$POSTGRES_PORT
  ./wait-for -t 15 $POSTGRES_HOST:$POSTGRES_PORT
fi

if [ ! -z "${REDIS_ADDRESS:-}" ]; then
  echo "wait for redis @" $REDIS_ADDRESS
  ./wait-for -t 15 $REDIS_ADDRESS
fi

if [ "$APP_NAME" = "tigerhall-kittens-server" ]; then
  ./tigerhall-kittens-server
else
  echo "unknown APP_NAME to start:" $APP_NAME
fi
