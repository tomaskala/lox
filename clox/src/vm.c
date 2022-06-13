#include <stdio.h>

#include "common.h"
#include "compiler.h"
#include "debug.h"
#include "vm.h"

VM vm;

static void
stack_reset()
{
  vm.stack_top = vm.stack;
}

void
vm_init()
{
  stack_reset();
}

static InterpretResult
run()
{
  #define READ_BYTE() (*vm.ip++)
  #define READ_CONSTANT() (vm.chunk->constants.values[READ_BYTE()])
  #define BINARY_OP(op) \
  do { \
    double b = vm_stack_pop(); \
    double a = vm_stack_pop(); \
    vm_stack_push(a op b); \
  } while (false)
  for (;;) {
    #ifdef DEBUG_TRACE_EXECUTION
    printf(" ");
    for (Value *slot = vm.stack; slot < vm.stack_top; ++slot) {
      printf("[ ");
      value_print(*slot);
      printf(" ]");
    }
    printf("\n");
    disassemble_instruction(vm.chunk, (size_t) (vm.ip - vm.chunk->code));
    #endif
    uint8_t instruction;
    switch (instruction = READ_BYTE()) {
    case OP_CONSTANT:
      {
        Value constant = READ_CONSTANT();
        vm_stack_push(constant);
        break;
      }
    case OP_ADD:
      BINARY_OP(+);
      break;
    case OP_SUBTRACT:
      BINARY_OP(-);
      break;
    case OP_MULTIPLY:
      BINARY_OP(*);
      break;
    case OP_DIVIDE:
      BINARY_OP(/);
      break;
    case OP_NEGATE:
      vm_stack_push(-vm_stack_pop());
      break;
    case OP_RETURN:
      value_print(vm_stack_pop());
      printf("\n");
      return INTERPRET_OK;
    }
  }
  #undef READ_BYTE
  #undef READ_CONSTANT
  #undef BINARY_OP
}

InterpretResult
vm_interpret(const char *source)
{
  Chunk chunk;
  chunk_init(&chunk);
  if (!compile(source, &chunk)) {
    chunk_free(&chunk);
    return INTERPRET_COMPILE_ERROR;
  }
  vm.chunk = &chunk;
  vm.ip = vm.chunk->code;
  InterpretResult result = run();
  chunk_free(&chunk);
  return result;
}

void
vm_stack_push(Value value)
{
  *vm.stack_top = value;
  vm.stack_top++;
}

Value
vm_stack_pop()
{
  vm.stack_top--;
  return *vm.stack_top;
}

void
vm_free()
{
}
