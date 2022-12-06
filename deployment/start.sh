#!/usr/bin/env sh

set -e

/opt/bin/tern migrate -m /opt/bin/migrations -c /opt/bin/migrations/tern.conf
exec /opt/bin/app
