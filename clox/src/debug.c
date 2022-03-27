#include <inttypes.h>
#include <stdio.h>

#include "debug.h"
#include "value.h"

void
disassemble_chunk(Chunk *chunk, const char *name)
{
  printf("== %s ==\n", name);
  for (size_t offset = 0; offset < chunk->count;)
    offset = disassemble_instruction(chunk, offset);
}

static size_t
constant_instruction(const char *name, Chunk *chunk, size_t offset)
{
  uint8_t constant = chunk->code[offset + 1];
  printf("%-16s %4" PRIu8 " '", name, constant);
  value_print(chunk->constants.values[constant]);
  printf("'\n");
  return offset + 2;
}

static size_t
simple_instruction(const char *name, size_t offset)
{
  printf("%s\n", name);
  return offset + 1;
}

size_t
disassemble_instruction(Chunk *chunk, size_t offset)
{
  printf("%04zu ", offset);
  uint8_t instruction = chunk->code[offset];
  switch (instruction) {
  case OP_CONSTANT:
    return constant_instruction("OP_CONSTANT", chunk, offset);
  case OP_RETURN:
    return simple_instruction("OP_RETURN", offset);
  default:
    printf("Unknown opcode %" PRIu8 "\n", instruction);
    return offset + 1;
  }
}
