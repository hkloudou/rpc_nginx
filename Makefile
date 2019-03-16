CLOSE_TAG = $(shell git describe --abbrev=0)
IMAGE_PREFIX=hkloudou
COMPONENT=rpc-nginx
IMAGE = $(IMAGE_PREFIX)/$(COMPONENT):$(CLOSE_TAG)
IMAGEWITHOUTTAG = $(IMAGE_PREFIX)/$(COMPONENT):latest
LDFLAGS = -X main._version_=$(CLOSE_TAG) -X main._branch_=$(GIT_BRANCH) -X main._commitId_=$(GIT_LAST_COMMIT) -X main._buildTime_=$(COMPILE_TIME) -X main._appName_=$(COMPONENT) -s -w
default: init

init:
	@git config --local user.name hkloudou
	@git config --local user.email hkloudou@gmail.com
	@git config --local user.signingkey 575A0CADC23D0A96
	@git config --local commit.gpgsign true
	@git config --local autotag.sign true
git:
	git autotag -commit 'rpc nginx with grpc' -tag=true -push=true
build:
	@make git
	mkdir -p docker/bin/ && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -a -o docker/bin/rpc-nginx ./server/
	docker build -t $(IMAGEWITHOUTTAG) .
deploy:
	@docker tag $(IMAGEWITHOUTTAG) $(IMAGE)
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
