VERSION = $(shell git autotag version)
VERSIONNILDOT = $(subst .,-,$(VERSION))
GIT_LAST_HASH = $(shell git rev-list --tags --max-count=1)
GIT_SHA = $(shell git rev-parse --short HEAD)
GIT_LAST_TAG = $(shell git describe --tags ${GIT_LAST_HASH})
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
GIT_LAST_COMMIT = $(shell git log --pretty=format:"%h" -1)
COMPILE_TIME = $(shell date +%s)

IMAGE_PREFIX=hkloudou
COMPONENT=rpc-nginx
WEBPROTO=nginx
IMAGE = $(IMAGE_PREFIX)/$(COMPONENT):$(VERSION)
IMAGEWITHOUTTAG = $(IMAGE_PREFIX)/$(COMPONENT):latest
LDFLAGS = -X main._version_=$(VERSION) -X main._branch_=$(GIT_BRANCH) -X main._commitId_=$(GIT_LAST_COMMIT) -X main._buildTime_=$(COMPILE_TIME) -X main._appName_=$(COMPONENT) -s -w
    
default: init
init:
	@git config --local user.name hkloudou
	@git config --local user.email hkloudou@gmail.com
	@git config --local user.signingkey 575A0CADC23D0A96
	@git config --local commit.gpgsign true
	@git config --local autotag.sign true
protoc:
	protoc -I $(WEBPROTO)/ $(WEBPROTO)/$(WEBPROTO).proto --go_out=plugins=grpc:$(WEBPROTO)
dist-clean:
	rm -rf docker/bin/
build-bin: dist-clean
	mkdir -p docker/bin/ && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o docker/bin/$(COMPONENT) ./server
build-debug-docker:
	cd docker && docker build --rm --build-arg VERSION=debug --build-arg COMPONENT=$(COMPONENT) -t $(IMAGEWITHOUTTAG) .
build-docker:
	cd docker && docker build --rm --build-arg VERSION=$(VERSIONNILDOT) --build-arg COMPONENT=$(COMPONENT) -t $(IMAGEWITHOUTTAG) .
deploy:
	-git autotag -commit 'auto commit' -t -i
	@make build-bin
	@make build-docker
	@make dist-clean
	@make deploy-docker
	@make fixversion
fixversion:
	echo "package $(WEBPROTO)">"$(WEBPROTO)/version.go"
	echo "" >> "$(WEBPROTO)/version.go"
	echo "//Version git tag version" >> "$(WEBPROTO)/version.go"
	echo "var Version = \"$(VERSION)\"" >>"$(WEBPROTO)/version.go"
	echo "var VersionNilDot = \"$(VERSIONNILDOT)\"" >>"$(WEBPROTO)/version.go"
	echo "var GRPCServerURL = \"${COMPONENT}-$(VERSIONNILDOT).grpc.apiatm.com\"" >>"$(WEBPROTO)/version.go"
	git autotag -commit 'auto commit' -t -f -p
deploy-docker:
	@docker tag $(IMAGEWITHOUTTAG) $(IMAGE)
	@docker push $(IMAGEWITHOUTTAG)
	@docker push $(IMAGE)
debug:
	@make down
	@make git
	@make build-bin
	@make build-debug-docker
	@make dist-clean
	@make up
up:
	@make down
	docker-compose up -d
u:
	@make down
	docker-compose up
down:
	docker-compose down -v --remove-orphans
t:
	echo  $(subst .,-,$(VERSION))
ssl:
	nginx-ssl-push -cert "${HOME}/.acme.sh/apiatm.com/fullchain.cer" -key "${HOME}/.acme.sh/apiatm.com/apiatm.com.key" -name "apiatm.com" -url "localhost:9000"
install:
	rm -rf $(GOPATH)/bin/nginx-ssl-push
	cd cmd/nginx-ssl-push && go install