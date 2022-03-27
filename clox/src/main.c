#include <stdio.h>

#include "chunk.h"
#include "debug.h"

int
main()
{
  #ifdef CLOX_DEBUG_H
  Chunk chunk;
  chunk_init(&chunk);
  size_t constant = chunk_add_constant(&chunk, 1.2);
  chunk_write(&chunk, OP_CONSTANT);
  chunk_write(&chunk, constant);
  chunk_write(&chunk, OP_RETURN);
  disassemble_chunk(&chunk, "test chunk");
  chunk_free(&chunk);
  #endif
  return 0;
}
