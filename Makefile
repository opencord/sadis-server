# Copyright 2018 Open Networking Foundation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

REGISTRY          ?=
REPOSITORY        ?=
DOCKER_BUILD_ARGS ?=
TAG               ?= $(shell cat ${MAKEFILE_DIR}/VERSION)
IMAGENAME         := ${REGISTRY}${REPOSITORY}sadis-server:${TAG}
SHELL             := /bin/bash

all: build push

build:
	docker build $(DOCKER_BUILD_ARGS) -t ${IMAGENAME} -f Dockerfile .

push:
	docker push ${IMAGENAME}