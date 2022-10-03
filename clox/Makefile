.POSIX:
.SUFFIXES:
CC      = cc
CFLAGS  = -std=c99 -pedantic -Wall -Wextra -Werror -Wno-unused-parameter
LDFLAGS =
LDLIBS  =

SRCS = \
	src/main.c \
	src/chunk.c \
	src/memory.c \
	src/debug.c \
	src/value.c \
	src/vm.c \
	src/compiler.c \
	src/scanner.c \
	src/object.c \
	src/table.c

DBGEXE    = dbg
DBGOBJS   = $(SRCS:.c=.dbg.o)
DBGCFLAGS = -Og -g -DDEBUG_PRINT_CODE -DDEBUG_TRACE_EXECUTION -DDEBUG_STRESS_GC -DDEBUG_LOG_GC

RELEXE    = lox
RELOBJS   = $(SRCS:.c=.o)
RELCFLAGS = -O3

.PHONY: all
all: $(RELEXE)

.PHONY: debug
debug: $(DBGEXE)

$(RELEXE): $(RELOBJS)
	$(CC) $(LDFLAGS) -o $(RELEXE) $(RELOBJS) $(LDLIBS)

$(DBGEXE): $(DBGOBJS)
	$(CC) $(LDFLAGS) -o $(DBGEXE) $(DBGOBJS) $(LDLIBS)

.PHONY: clean
clean:
	rm -f $(RELEXE) $(RELOBJS) $(DBGEXE) $(DBGOBJS)

.SUFFIXES: .c .o
.c.o:
	$(CC) $(CFLAGS) $(RELCFLAGS) -o $@ -c $<

.SUFFIXES: .c .dbg.o
.c.dbg.o:
	$(CC) $(CFLAGS) $(DBGCFLAGS) -o $@ -c $<
