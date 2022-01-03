DC = docker-compose
DC_WEB = $(DC) exec web
DC_DB = $(DC) exec db

package_name:
	@$(eval PACKAGE_NAME := $(shell read -p "add packages: " NAME; echo $$NAME))

go_get: package_name
	${DC_WEB} go get ${PACKAGE_NAME}

go_install: package_name
	${DC_WEB} go install ${PACKAGE_NAME}

mod_tidy:
	${DC_WEB} go mod tidy

db_attach:
	${DC_DB} bash