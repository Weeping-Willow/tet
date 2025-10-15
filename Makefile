.PHONY: deps goose

LOCAL_BIN := $(CURDIR)/bin
GOOSE_VERSION := v3.26.0

deps: deps-goose

deps-goose:
ifeq ("$(wildcard $(LOCAL_BIN)/goose)","")
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
endif

