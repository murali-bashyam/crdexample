
ifeq ($(GO_PROJECT),)
$(error the variable GO_PROJECT must be set prior to including golang.mk)
endif

# Optional. These are subdirs that we look for all go files to test, vet, and fmt
GO_SUBDIRS ?= cmd pkg

# Optional directories (relative to CURDIR)
GO_PKG_DIR ?= $(WORK_DIR)/pkg

GO_SUPPORTED_VERSIONS ?= 1.15

GO_PACKAGES := $(foreach t,$(GO_SUBDIRS),$(GO_PROJECT)/$(t)/...)

GO := go
GOHOST := GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) go
GO_VERSION := $(shell $(GO) version | sed -ne 's/[^0-9]*\(\([0-9]\.\)\{0,4\}[0-9][^.]\).*/\1/p')

GOLINT := $(TOOLS_HOST_DIR)/golint

# we use a consistent version of gofmt even while running different go compilers.
# see https://github.com/golang/go/issues/26397 for more details
GOFMT_VERSION := 1.11
ifneq ($(findstring $(GOFMT_VERSION),$(GO_VERSION)),)
GOFMT := $(shell which gofmt)
else
GOFMT := $(TOOLS_HOST_DIR)/gofmt$(GOFMT_VERSION)
endif

GO_OUT_DIR := $(abspath $(OUTPUT_DIR)/bin/$(PLATFORM))

.PHONY: go.build
go.build:
	@echo === go build $(PLATFORM)
	$(foreach p,$(GO_STATIC_PACKAGES),@CGO_ENABLED=0 $(GO) build -v -i -o $(GO_OUT_DIR)/$(lastword $(subst /, ,$(p)))$(GO_OUT_EXT) $(GO_STATIC_FLAGS) $(p)${\n})
	$(foreach p,$(GO_TEST_PACKAGES),@CGO_ENABLED=0 $(GO) test -v -i -c -o $(GO_TEST_OUTPUT)/$(lastword $(subst /, ,$(p)))$(GO_OUT_EXT) $(GO_STATIC_FLAGS) $(p)${\n})

.PHONY: go.lint
go.lint: 
	@echo === go lint
	@$(GOLINT) -set_exit_status=true $(GO_PACKAGES) 

.PHONY: go.vet
go.vet:
	@echo === go vet
	@CGO_ENABLED=0 $(GOHOST) vet $(GO_COMMON_FLAGS) $(GO_PACKAGES) 

.PHONY: go.fmt
go.fmt: $(GOFMT)
	@gofmt_out=$$($(GOFMT) -s -d -e $(GO_SUBDIRS) 2>&1) && [ -z "$${gofmt_out}" ] || (echo "$${gofmt_out}" 1>&2; exit 1)
go.validate: go.vet go.fmt

.PHONY: go.mod.update
go.mod.update:
	@echo === updating modules
	@$(GOHOST) get -u ./...

.PHONY: go.mod.check
go.mod.check:
	@echo === ensuring modules are tidied
	@$(GOHOST) mod tidy

.PHONY: go.mod.clean
go.mod.clean:
	@echo === cleaning modules cache
	@sudo rm -fr $(WORK_DIR)/cross_pkg
	@$(GOHOST) clean -modcache
