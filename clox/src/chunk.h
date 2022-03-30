#ifndef CLOX_CHUNK_H
#define CLOX_CHUNK_H

#include "common.h"
#include "value.h"

typedef enum {
  OP_CONSTANT,
  OP_ADD,
  OP_SUBTRACT,
  OP_MULTIPLY,
  OP_DIVIDE,
  OP_NEGATE,
  OP_RETURN,
} OpCode;

typedef struct {
  size_t count;
  size_t capacity;
  uint8_t *code;
  int *lines;
  ValueArray constants;
} Chunk;

void
chunk_init(Chunk *chunk);

void
chunk_write(Chunk *chunk, uint8_t byte, int line);

size_t
chunk_add_constant(Chunk *chunk, Value value);

void
chunk_free(Chunk *chunk);

#endif
