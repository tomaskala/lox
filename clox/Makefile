.POSIX:
.SUFFIXES:
CC      = cc
CFLAGS  = -std=c99 -pedantic -Wall -Wextra -Werror -Wno-unused-parameter
LDFLAGS =
LDLIBS  =

sources = \
	src/main.c \
	src/chunk.c \
	src/memory.c \
	src/debug.c \
	src/value.c \
	src/vm.c \
	src/compiler.c \
	src/scanner.c
headers = \
	src/common.h \
	src/chunk.h \
	src/memory.h \
	src/debug.h \
	src/value.h \
	src/vm.h \
	src/compiler.h \
	src/scanner.h \
	src/object.h
objects = $(sources:.c=.o)

all: clox

clox: $(objects)
	$(CC) -O3 $(LDFLAGS) -o $@ $(objects) $(LDLIBS)

debug: $(objects)
	$(CC) -g -Og $(LDFLAGS) -o $@ $(objects) $(LDLIBS)

.PHONY: clean
clean:
	rm -f clox debug $(objects)

.SUFFIXES: .c .o
.c.o:
	$(CC) -c $(CFLAGS) -o $@ $<
