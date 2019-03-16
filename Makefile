CLOSE_TAG = $(shell git describe --abbrev=0)
IMAGE_PREFIX=hkloudou
COMPONENT=rpc-nginx
IMAGE = $(IMAGE_PREFIX)/$(COMPONENT):$(CLOSE_TAG)
IMAGEWITHOUTTAG = $(IMAGE_PREFIX)/$(COMPONENT):latest
default: init

init:
	@git config --local user.name hkloudou
	@git config --local user.email hkloudou@gmail.com
	@git config --local user.signingkey 575A0CADC23D0A96
	@git config --local commit.gpgsign true
	@git config --local autotag.sign true
git:
	git autotag -commit 'nginx alpine 1.15.9' -tag=true -push=true
build:
	@make git
	docker build -t $(IMAGEWITHOUTTAG) .
tag:
	@docker tag $(IMAGEWITHOUTTAG) $(IMAGE)
deploy:
	@docker push $(IMAGEWITHOUTTAG)
	@docker push $(IMAGE)
up:
	@make down
	docker-compose up -d
u:
	@make down
	docker-compose up
down:
	docker-compose down -v --remove-orphans
t:
	curl -H "Host: whoami.local" localhost
protoc:
	protoc -I proto/ proto/nginx.proto --go_out=plugins=grpc:proto
