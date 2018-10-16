REGISTRY := gcr.io/kouzoh-p-manoj-gubba
SERVICE_NAME := item-service
VERSION := latest

IMAGE := $(REGISTRY)/$(SERVICE_NAME):$(VERSION)

.PHONY: dep
dep:
	dep ensure -v

.PHONY: build
build:
	GOOS=linux go build -o item-service

.PHONY: docker/build
docker/build:
	docker build -t $(IMAGE) .

.PHONY: docker/push
docker/push:
	docker push $(IMAGE)

.PHONY: docker/network
docker/network:
	docker network create item-service

.PHONY: run/mysql
run/mysql:
	docker run -d --rm --name mysql \
		--network item-service \
		-v $(PWD)/mysql-data:/var/lib/mysql \
		-p 3306:3306 \
		-e MYSQL_ROOT_PASSWORD=allforone \
		mysql:5.6

.PHONY: run/item-service
run/item-service:
	docker run --rm --name item-service \
	--network item-service \
	  -p 8080:8080 \
		-e MYSQL_HOST=mysql \
		-e MYSQL_PASS=allforone \
		$(IMAGE)
