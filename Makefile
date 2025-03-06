# Change these variables as necessary.
main_package_path = ./cmd
binary_name = kata_test

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:	
	@echo 'Usage:'	
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
    @echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	@test -z "$(shell git status --porcelain)"


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #


## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## build: build the application
.PHONY: build
build:
    # Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
		go build -o=/tmp/bin/${binary_name} ${main_package_path}


.PHONY: docker-build
docker-build:
	sudo docker compose --profile local-environment up  --build  

.PHONY: docker-stop
docker-stop:
	sudo docker compose --profile local-environment down --volumes

## run: run the  application
.PHONY: run
run: build
	/tmp/bin/${binary_name}

.PHONY: lint
lint:
	golangci-lint run
    

# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to the remote Git repository
.PHONY: push
push: confirm audit no-dirty
	git push

## production/deploy: deploy the application to production
.PHONY: production/deploy
production/deploy: confirm audit no-dirty
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=/tmp/bin/linux_amd64/${binary_name} ${main_package_path}
		upx -5 /tmp/bin/linux_amd64/${binary_name}
    # Include additional deployment steps here...