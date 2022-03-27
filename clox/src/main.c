#include <stdio.h>

#include "chunk.h"
#include "debug.h"

int
main()
{
  Chunk chunk;
  chunk_init(&chunk);
  chunk_write(&chunk, OP_RETURN);
  #ifdef CLOX_DEBUG_H
  disassemble_chunk(&chunk, "test chunk");
  #endif
  chunk_free(&chunk);
  return 0;
}
