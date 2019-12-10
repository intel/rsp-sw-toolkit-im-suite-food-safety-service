# Apache v2 license
#  Copyright (C) <2019> Intel Corporation
#
#  SPDX-License-Identifier: Apache-2.0
#

PROJECT_NAME ?= food-safety-service

# The default flags to use when calling submakes
GNUMAKEFLAGS = --no-print-directory

GIT_SHA = $(shell git rev-parse HEAD)

GO_FILES = $(shell find . -type f -name '*.go')
RES_FILES = $(shell find res/ -type f)

PROXY_ARGS =	--build-arg http_proxy=$(http_proxy) \
				--build-arg https_proxy=$(https_proxy) \
				--build-arg no_proxy=$(no_proxy) \
				--build-arg HTTP_PROXY=$(HTTP_PROXY) \
				--build-arg HTTPS_PROXY=$(HTTPS_PROXY) \
				--build-arg NO_PROXY=$(NO_PROXY)

EXTRA_BUILD_ARGS ?=

touch_target_file = mkdir -p $(@D) && touch $@

trap_ctrl_c = trap 'exit 0' INT;

compose = docker-compose

log = docker-compose logs $1 $2 2>&1

.PHONY: build clean iterate iterate-d tail start stop rm deploy kill down fmt ps

default: build

build: build/docker

clean:
	rm -rf build/*

build/docker:
	docker build --rm \
		--build-arg GIT_TOKEN=$(GIT_TOKEN) \
		$(PROXY_ARGS) \
		$(EXTRA_BUILD_ARGS) \
		-f Dockerfile_dev \
		--label "git_sha=$(GIT_SHA)" \
		-t rsp/$(PROJECT_NAME):dev \
		.
	@$(touch_target_file)

iterate: build up

iterate-d: build up-d
	$(trap_ctrl_c) $(MAKE) tail

restart:
	$(compose) restart $(args)

kill:
	$(compose) kill $(args)

tail:
	$(trap_ctrl_c) $(call log,-f --tail=10, $(args))

down:
	docker stack rm Food-Safety-Stack

up: build
	$(compose) up --remove-orphans $(args)

up-d: build
	$(MAKE) up args="-d $(args)"

deploy: 
	docker stack deploy \
		--with-registry-auth \
		--compose-file docker-compose.yml \
		Food-Safety-Stack

fmt:
	go fmt ./...

ps:
	$(compose) ps
