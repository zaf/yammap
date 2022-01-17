#
#	Copyright (C) 2022, Lefteris Zafiris <zaf@fastmail.com>
#	This program is free software, distributed under the terms of
#	the GNU GPL v3 License. See the LICENSE file
#	at the top of the source tree.
#

# Run Go tests in all supported platforms using Docker and Qemu. 

SRC_LOCAL:=$(shell pwd)
SRC:="/usr/src/"

CACHE_LOCAL:="$${HOME}/.cache/go-build/"
CACHE:="/root/.cache/go-build/"


test: ## Run tests for all supported platforms.
test: test_linux_amd64 test_linux_386 test_linux_arm64 test_linux_arm

test_linux_amd64:
	@echo "== Testing Linux/amd64 =="
	@docker pull --platform=linux/amd64 golang
	-@docker run --rm -w ${SRC} --platform=linux/amd64 --mount type=bind,source="${SRC_LOCAL}",target=${SRC} --mount type=bind,source=${CACHE_LOCAL},target=${CACHE} golang go test

test_linux_386:
	@echo "== Testing Linux/386 =="
	@docker pull --platform=linux/386 golang
	-@docker run --rm -w ${SRC} --platform=linux/386 --mount type=bind,source="${SRC_LOCAL}",target=${SRC} --mount type=bind,source=${CACHE_LOCAL},target=${CACHE} golang go test

test_linux_arm64:
	@echo "== Testing Linux/arm64 =="
	@docker pull --platform=linux/arm64 golang
	-@docker run --rm -w ${SRC} --platform=linux/arm64 --mount type=bind,source="${SRC_LOCAL}",target=${SRC} --mount type=bind,source=${CACHE_LOCAL},target=${CACHE} golang go test

test_linux_arm:
	@echo "== Testingh Linux/arm =="
	@docker pull --platform=linux/arm/v7 golang
	-@docker run --rm -w ${SRC} --platform=linux/arm/v7 --mount type=bind,source="${SRC_LOCAL}",target=${SRC} --mount type=bind,source=${CACHE_LOCAL},target=${CACHE} golang go test

help: ## Show this help.
	@egrep '^(.+)\:\ ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#'

.PHONY: test test_linux_386 test_linux_amd64 test_linux_arm test_linux_arm64 help