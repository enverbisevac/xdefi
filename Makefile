help:           ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

image: ## build docker image
	docker build -t sunmicrosystem/xdefitask:latest .

push:  ## push docker image to docker hub
	docker push sunmicrosystem/xdefitask:latest

start: ## start docker environment
	docker-compose -p xdefi up -d

stop:  ## stop docker environment
	docker-compose down

build: ## build server for local run
	CGO_ENABLED=0 go build -ldflags="-w -s" -a -installsuffix cgo -o server ./main.go
