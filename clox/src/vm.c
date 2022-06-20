#include <stdarg.h>
#include <stdio.h>
#include <string.h>

#include "common.h"
#include "compiler.h"
#include "memory.h"
#include "object.h"
#include "vm.h"

#ifdef DEBUG_TRACE_EXECUTION
#include "debug.h"
#endif

VM vm;

static void
vm_stack_reset()
{
  vm.stack_top = vm.stack;
}

static void
runtime_error(const char *format, ...)
{
  va_list args;
  va_start(args, format);
  vfprintf(stderr, format, args);
  va_end(args);
  fputs("\n", stderr);
  size_t instruction = vm.ip - vm.chunk->code - 1;
  int line = vm.chunk->lines[instruction];
  fprintf(stderr, "[line %d] in script\n", line);
  vm_stack_reset();
}

void
vm_init()
{
  vm_stack_reset();
}

void
vm_free()
{
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

static Value
vm_stack_peek(size_t distance)
{
  return vm.stack_top[-1 - distance];
}

static bool
is_falsey(Value value) {
  return IS_NIL(value) || (IS_BOOL(value) && !AS_BOOL(value));
}

static void
concatenate()
{
  ObjString *b = AS_STRING(vm_stack_pop());
  ObjString *a = AS_STRING(vm_stack_pop());
  size_t length = a->length + b->length;
  char *chars = ALLOCATE(char, length + 1);
  memcpy(chars, a->chars, a->length);
  memcpy(chars + a->length, b->chars, b->length);
  chars[length] = '\0';
  ObjString *result = take_string(chars, length);
  vm_stack_push(OBJ_VAL(result));
}

static InterpretResult
run()
{
  #define READ_BYTE() (*vm.ip++)
  #define READ_CONSTANT() (vm.chunk->constants.values[READ_BYTE()])
  #define BINARY_OP(value_type, op) \
  do { \
    if (!IS_NUMBER(vm_stack_peek(0)) || !IS_NUMBER(vm_stack_peek(1))) { \
      runtime_error("Operands must be numbers."); \
      return INTERPRET_RUNTIME_ERROR; \
    } \
    double b = AS_NUMBER(vm_stack_pop()); \
    double a = AS_NUMBER(vm_stack_pop()); \
    vm_stack_push(value_type(a op b)); \
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
    case OP_CONSTANT: {
      Value constant = READ_CONSTANT();
      vm_stack_push(constant);
      break;
    }
    case OP_NIL:
      vm_stack_push(NIL_VAL);
      break;
    case OP_TRUE:
      vm_stack_push(BOOL_VAL(true));
      break;
    case OP_FALSE:
      vm_stack_push(BOOL_VAL(false));
      break;
    case OP_EQUAL: {
      Value b = vm_stack_pop();
      Value a = vm_stack_pop();
      vm_stack_push(BOOL_VAL(values_equal(a, b)));
      break;
    }
    case OP_GREATER:
      BINARY_OP(BOOL_VAL, >);
      break;
    case OP_LESS:
      BINARY_OP(BOOL_VAL, <);
      break;
    case OP_ADD: {
      if (IS_STRING(vm_stack_peek(0)) && IS_STRING(vm_stack_peek(1))) {
        concatenate();
      } else if (IS_NUMBER(vm_stack_peek(0)) && IS_NUMBER(vm_stack_peek(1))) {
        double b = AS_NUMBER(vm_stack_pop());
        double a = AS_NUMBER(vm_stack_pop());
        vm_stack_push(NUMBER_VAL(a + b));
      } else {
        runtime_error("Operands must be two numbers or two strings.");
        return INTERPRET_RUNTIME_ERROR;
      }
      break;
    }
    case OP_SUBTRACT:
      BINARY_OP(NUMBER_VAL, -);
      break;
    case OP_MULTIPLY:
      BINARY_OP(NUMBER_VAL, *);
      break;
    case OP_DIVIDE:
      BINARY_OP(NUMBER_VAL, /);
      break;
    case OP_NOT:
      vm_stack_push(BOOL_VAL(is_falsey(vm_stack_pop())));
      break;
    case OP_NEGATE:
      if (!IS_NUMBER(vm_stack_peek(0))) {
        runtime_error("Operand must be a number.");
        return INTERPRET_RUNTIME_ERROR;
      }
      vm_stack_push(NUMBER_VAL(-AS_NUMBER(vm_stack_pop())));
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
