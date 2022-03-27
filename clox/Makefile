.POSIX:
.SUFFIXES:
CC      = cc
CFLAGS  = -std=c99 -pedantic -Wall -Wextra -Werror
LDFLAGS =
LDLIBS  =

sources = src/main.c
headers =
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
