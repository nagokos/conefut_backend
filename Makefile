DC = docker-compose
DC_WEB = $(DC) exec web
DC_DB = $(DC) exec db

package_add:
	@$(eval PACKAGE_NAME := $(shell read -p "add packages: " NAME; echo $$NAME))

package_name:
	@$(eval PACKAGE_NAME := $(shell read -p "packages name: " NAME; echo $$NAME))

go_run: package_name
	${DC_WEB} go run ${PACKAGE_NAME}

go_get: package_add
	${DC_WEB} go get -d ${PACKAGE_NAME}

go_install: package_add
	${DC_WEB} go install ${PACKAGE_NAME}

start_server:
	${DC_WEB} go run graph/server/server.go

mod_tidy:
	${DC_WEB} go mod tidy

generate:
	${DC_WEB} go generate ./...

db_attach:
	${DC_DB} bash

