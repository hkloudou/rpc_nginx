CLOSE_TAG = $(shell git describe --abbrev=0)

GIT_LAST_HASH = $(shell git rev-list --tags --max-count=1)
GIT_SHA = $(shell git rev-parse --short HEAD)
GIT_LAST_TAG = $(shell git describe --tags ${GIT_LAST_HASH})
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
GIT_LAST_COMMIT = $(shell git log --pretty=format:"%h" -1)
COMPILE_TIME = $(shell date +%s)

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
	@make build-bin
build-bin:
	mkdir -p docker/bin/ && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -a -o docker/bin/rpc-nginx ./server/
	cd docker && docker build -t $(IMAGEWITHOUTTAG) .
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
	nginx-ssl-push -cert "${HOME}/.acme.sh/apiatm.com/fullchain.cer" -key "${HOME}/.acme.sh/apiatm.com/apiatm.com.key" -name "apiatm.com" -url "localhost:9000"
	nginx-ssl-push -cert "${HOME}/.acme.sh/youziku.com/fullchain.cer" -key "${HOME}/.acme.sh/youziku.com/youziku.com.key" -name "youziku.com" -url "localhost:9000"
protoc:
	protoc -I proto/ proto/nginx.proto --go_out=plugins=grpc:proto
install:
	cd cmd/nginx-ssl-push && go install