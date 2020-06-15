GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get ./...
SERVICE=workscheduler
DIST=dist
BINARY=$(DIST)/$(SERVICE)
BINARY_LINUX=$(BINARY)_linux
DOCKER_REPOSITORY=tidepool/$(SERVICE)

all: test build
ci:	test docker-build docker-push-ci
build:	deps
		$(GOBUILD) -o $(BINARY) -v ./server
linux-build:	deps
		env CGO_ENABLED=1 GOOS="linux" GOARCH="amd64" $(GOBUILD) -v -a -tags musl -o $(BINARY_LINUX) ./server
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY)
start:	build
		./$(BINARY)
deps:
		$(GOGET)
docker-login:
		@echo "$(DOCKER_PASSWORD)" | docker login --username "$(DOCKER_USERNAME)" --password-stdin
docker-build:	linux-build
		docker build -t $(SERVICE) .
docker-push-ci:	docker-login
ifdef TRAVIS_BRANCH
ifdef TRAVIS_COMMIT
ifdef TRAVIS_PULL_REQUEST_BRANCH
	docker tag $(SERVICE) $(DOCKER_REPOSITORY):PR-$(subst /,-,$(TRAVIS_BRANCH))-$(TRAVIS_COMMIT)
	docker push $(DOCKER_REPOSITORY):PR-$(subst /,-,$(TRAVIS_BRANCH))-$(TRAVIS_COMMIT)
else
	docker tag $(SERVICE) $(DOCKER_REPOSITORY):$(subst /,-,$(TRAVIS_BRANCH))-$(TRAVIS_COMMIT)
	docker tag $(SERVICE) $(DOCKER_REPOSITORY):$(subst /,-,$(TRAVIS_BRANCH))-latest
	docker push $(DOCKER_REPOSITORY):$(subst /,-,$(TRAVIS_BRANCH))-$(TRAVIS_COMMIT)
	docker push $(DOCKER_REPOSITORY):$(subst /,-,$(TRAVIS_BRANCH))-latest
endif
endif
endif
