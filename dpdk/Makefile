# DPDK TCP Forwarder Makefile
#
# Adjust DPDK_PATH and RTE_SDK-style vars to match your environment.
# Tested against DPDK 21.11 LTS / 22.x.

DPDK_PATH ?= /usr/local/share/dpdk
PKGCONF   ?= pkg-config

ifeq ($(shell $(PKGCONF) --exists libdpdk && echo yes),yes)
CFLAGS  += $(shell $(PKGCONF) --cflags libdpdk)
LDFLAGS += $(shell $(PKGCONF) --libs libdpdk)
else
# Fall back to manual flags. Set RTE_INCLUDE and RTE_LIB accordingly.
RTE_INCLUDE ?= /usr/local/include
RTE_LIB     ?= /usr/local/lib
CFLAGS  += -I$(RTE_INCLUDE) -march=native
LDFLAGS += -L$(RTE_LIB) -ldpdk -lpthread -ldl -lm
endif

CFLAGS  += -O2 -g -Wall -Wextra -std=gnu11
LDFLAGS += -export-dynamic

BIN := build/tcp_forwarder

all: $(BIN)

build:
	mkdir -p build

$(BIN): tcp_forwarder.c | build
	$(CC) $(CFLAGS) $< -o $@ $(LDFLAGS)

clean:
	rm -rf build

.PHONY: all clean