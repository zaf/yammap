#
#	Copyright (C) 2022, Lefteris Zafiris <zaf@fastmail.com>
#	This program is free software, distributed under the terms of
#	the GNU GPL v3 License. See the LICENSE file
#	at the top of the source tree.
#

# Run Go tests in all supported platforms using Docker and Qemu. 

SRC_LOCAL:=$(shell pwd)
SRC:="/usr/src/"

CACHE_LOCAL:=$(shell go env GOCACHE)
CACHE:="/root/.cache/go-build/"

export DOCKER_BUILDKIT=1

define run-test =
	@echo "== Testing $(1) =="
	@docker pull --platform=$(1) golang > /dev/null
	-@docker run --rm -w ${SRC} --platform=$(1) --mount type=bind,source="${SRC_LOCAL}",target=${SRC} --mount type=bind,source=${CACHE_LOCAL},target=${CACHE} golang go test
endef

test: ## Run tests for all supported platforms.
test: test_linux_amd64 test_linux_386 test_linux_arm64 test_linux_arm test_linux_mips64le

test_linux_amd64: ## Run tests for linux/amd64
	$(call run-test,linux/amd64)

test_linux_386: ## Run tests for linux/386
	$(call run-test,linux/386)

test_linux_arm64: ## Run tests for linux/arm64
	$(call run-test,linux/arm64)

test_linux_arm: ## Run tests for linux/arm
	$(call run-test,linux/arm/v7)

test_linux_mips64le: ## Run tests for linux/mips64le
	$(call run-test,linux/mips64le)

help: ## Show this help.
	@egrep '^(.+)\:\ ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#'

.PHONY: test test_linux_386 test_linux_amd64 test_linux_arm test_linux_arm64 test_linux_mips64le help
