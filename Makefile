all: build test

include tasks/build.mk
include tasks/client.mk
include tasks/db.mk
include tasks/dev.mk
