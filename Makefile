MAKEFLAGS += -r --warn-undefined-variables
SHELL := /bin/bash
.SHELLFLAGS := -o pipefail -euc
.DEFAULT_GOAL := help

PROJECT = logrus-kafka-hook
export IMPORT_PATH = github.com/kenjones-cisco/logrus-kafka-hook

# Windows environment?
CYG_CHECK := $(shell hash cygpath 2>/dev/null && echo 1)
# WSL environment?
WSL_CHECK := $(shell grep -qE "(Microsoft|WSL)" /proc/version 2>/dev/null && echo 1)
ifeq ($(CYG_CHECK),1)
	VBOX_CHECK := $(shell hash VBoxManage 2>/dev/null && echo 1)

	# Docker Toolbox (pre-Windows 10)
	ifeq ($(VBOX_CHECK),1)
		ROOT := /${PROJECT}
	else
		# Docker Windows
		ROOT := $(shell cygpath -m -a "$(shell pwd)")
	endif
else ifeq ($(WSL_CHECK),1)
	# requires drives shared in Docker for Windows settings
	ROOT := $(strip $(shell cmd.exe /c cd | sed -e 's/\\/\//g'))
else
	# all non-windows environments
	ROOT := $(shell pwd)
endif

DEV_IMAGE := ${PROJECT}_dev

DOCKERRUN := docker run --rm \
	-v ${ROOT}/vendor:/go/src \
	-v ${ROOT}:/${PROJECT}/src/${IMPORT_PATH} \
	-w /${PROJECT}/src/${IMPORT_PATH} \
	${DEV_IMAGE}

DOCKERNOVENDOR := docker run --rm -i \
	-e IMPORT_PATH="${IMPORT_PATH}" \
	-v ${ROOT}:/${PROJECT}/src/${IMPORT_PATH} \
	-w /${PROJECT}/src/${IMPORT_PATH} \
	${DEV_IMAGE}


.PHONY: help clean veryclean vendor dep-* format check test cover

## display this help message
help:
	@echo 'Management commands for logrus-kafka-hook:'
	@echo
	@echo 'Usage:'
	@echo '  ## Develop / Test Commands'
	@echo '    vendor          Install dependencies using glide if glide.yaml changed.'
	@echo '    dep-update      Update dependencies using glide.'
	@echo '    dep-add         Add new dependencies to glide and install.'
	@echo '    format          Run code formatter.'
	@echo '    check           Run static code analysis (lint).'
	@echo '    test            Run tests on project.'
	@echo '    cover           Run tests and capture code coverage metrics on project.'
	@echo '    clean           Clean the directory tree of produced artifacts.'
	@echo '    veryclean       Same as clean but also removes cached dependencies.'
	@echo

## prefix before other make targets to run in your local dev environment
local: | quiet
	@$(eval DOCKERRUN= )
	@$(eval DOCKERNOVENDOR= )
	@mkdir -p tmp
	@date > tmp/dev_image_id
quiet: # this is silly but shuts up 'Nothing to be done for `local`'
	@:

## Clean the directory tree of produced artifacts.
clean:
	@rm -rf cover *.out *.xml

## Same as clean but also removes cached dependencies.
veryclean: clean
	@rm -rf tmp .glide vendor

## builds the dev container
prepare: tmp/dev_image_id
tmp/dev_image_id: Dockerfile.dev
	@mkdir -p tmp
	@docker rmi -f ${DEV_IMAGE} > /dev/null 2>&1 || true
	@echo "## Building dev container"
	@docker build --quiet -t ${DEV_IMAGE} -f Dockerfile.dev .
	@docker inspect -f "{{ .ID }}" ${DEV_IMAGE} > tmp/dev_image_id

# ----------------------------------------------
# dependencies

## Install dependencies using glide if glide.yaml changed.
vendor: tmp/glide-installed
tmp/glide-installed: tmp/dev_image_id glide.yaml
	@mkdir -p vendor
	${DOCKERRUN} glide install --skip-test --strip-vendor
	@date > tmp/glide-installed
	@chmod 644 glide.lock || :

## Update dependencies using glide.
dep-update: prepare
	${DOCKERRUN} glide up --skip-test --strip-vendor
	@chmod 644 glide.lock || :

# usage DEP=github.com/owner/package make dep-add
## Add new dependencies to glide and install.
dep-add: prepare
ifeq ($(strip $(DEP)),)
	$(error "No dependency provided. Expected: DEP=<go import path>")
endif
	${DOCKERNOVENDOR} glide get --skip-test --strip-vendor ${DEP}
	@chmod 644 glide.lock || :

# ----------------------------------------------
# develop and test

## print environment info about this dev environment
debug:
	@echo IMPORT_PATH="$(IMPORT_PATH)"
	@echo ROOT="$(ROOT)"
	@echo
	@echo docker commands run as:
	@echo "$(DOCKERRUN)"

## Run code formatter.
format: tmp/glide-installed
	${DOCKERNOVENDOR} bash ./scripts/format.sh

## Run static code analysis (lint).
check: format
	${DOCKERNOVENDOR} bash ./scripts/check.sh

## Run tests on project.
test: check
	${DOCKERRUN} bash ./scripts/test.sh

## Run tests and capture code coverage metrics on project.
cover: check
	@rm -rf cover/
	@mkdir -p cover
	${DOCKERRUN} bash ./scripts/cover.sh
	@chmod 644 cover/coverage.html
