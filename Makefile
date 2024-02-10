SHELL=/bin/bash -o pipefail
export SQLITE_DB?=./app.db

include config.mk

run:
	go run cmd/foo/main.go

sql:
	sqlite3 ${SQLITE_DB}

${SQLITE_DB}:
	cat internal/db/schema.sql | sqlite3 $@

gen:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

reset: dist-clean ${SQLITE_DB} gen

dist-clean:
	rm -f ${SQLITE_DB}
