SHELL=/bin/bash -o pipefail

include config.mk

run:
	SQLITE_DB=./app.db go run cmd/megawork/main.go

sql:
	sqlite3 app.db

app.db:
	cat internal/db/schema.sql | sqlite3 $@

gen:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

resetdb: dist-clean app.db gen

dist-clean:
	rm -f app.db
