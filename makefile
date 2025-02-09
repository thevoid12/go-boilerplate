
include .env
export $(shell sed 's/=.*//' .env)

# Shell and Make defaults
SHELL				:= /bin/bash

migrate-up:
	@echo "**************************** migration up ***************************************"
	@command="goose -dir migrations postgres \"user=${PG_USER} password=${PG_PASSWORD} dbname=${PG_DB} sslmode=${PG_SSLMODE}\" up"; \
	echo $$command; \
	result=$$(eval $$command); \
	echo "$$result"
	@echo "******************************************************************************"
migrate-down:
	@echo "**************************** migration down ***************************************"
	@command="goose -dir migrations postgres \"user=${PG_USER} password=${PG_PASSWORD} dbname=${PG_DB} sslmode=${PG_SSLMODE}\" down"; \
	echo $$command; \
	result=$$(eval $$command); \
	echo "$$result"
	@echo "******************************************************************************"
