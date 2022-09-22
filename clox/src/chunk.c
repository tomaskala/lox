#include <stdlib.h>

#include "chunk.h"
#include "memory.h"
#include "vm.h"

void
chunk_init(Chunk *chunk)
{
  chunk->count = 0;
  chunk->capacity = 0;
  chunk->code = NULL;
  chunk->lines = NULL;
  value_array_init(&chunk->constants);
}

void
chunk_write(Chunk *chunk, uint8_t byte, int line)
{
  if (chunk->capacity < chunk->count + 1) {
    int old_capacity = chunk->capacity;
    chunk->capacity = GROW_CAPACITY(old_capacity);
    chunk->code = GROW_ARRAY(uint8_t, chunk->code, old_capacity,
                             chunk->capacity);
    chunk->lines = GROW_ARRAY(int, chunk->lines, old_capacity,
                              chunk->capacity);
  }
  chunk->code[chunk->count] = byte;
  chunk->lines[chunk->count] = line;
  chunk->count++;
}

int
chunk_add_constant(Chunk *chunk, Value value)
{
  vm_stack_push(value);
  value_array_write(&chunk->constants, value);
  vm_stack_pop();
  return chunk->constants.count - 1;
}

void
chunk_free(Chunk *chunk)
{
  FREE_ARRAY(uint8_t, chunk->code, chunk->capacity);
  FREE_ARRAY(int, chunk->lines, chunk->capacity);
  value_array_free(&chunk->constants);
  chunk_init(chunk);
}
