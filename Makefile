BUILD_DIR := $(shell echo `pwd`/build)
UNAME_S := $(shell uname -s)

DYNAMIC_BUILD_FLAGS := -ldflags "-X work/gobox/utils.mode=DYNAMIC"
STATIC_BUILD_FLAGS := -ldflags "-linkmode external -extldflags -static -X work/gobox/utils.mode=STATIC"

ifeq ($(UNAME_S),Linux)
	BUILD_FLAGS = $(STATIC_BUILD_FLAGS)
else
	BUILD_FLAGS = $(DYNAMIC_BUILD_FLAGS)
endif

GO_INSTALL_ENV := GOBIN=$(BUILD_DIR) GOCACHE=$(BUILD_DIR)/cache GOGC=100 GOFLAGS="-p=2"
GO_RUN_ENV := GOBIN=$(BUILD_DIR) GOCACHE=$(BUILD_DIR)/cache GOGC=100 GOFLAGS="-p=2"

init: | $(BUILD_DIR)

gobox: init
	$(GO_RUN_ENV) go run $(BUILD_FLAGS) work/gobox

gobox-install:
	$(GO_RUN_ENV) go install $(BUILD_FLAGS) work/gobox
