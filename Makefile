all: build test

include tasks/client.mk
include tasks/build.mk
include tasks/dev.mk
