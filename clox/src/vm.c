#include <stdarg.h>
#include <stdio.h>
#include <string.h>
#include <time.h>

#include "common.h"
#include "compiler.h"
#include "memory.h"
#include "object.h"
#include "vm.h"

#ifdef DEBUG_TRACE_EXECUTION
#include "debug.h"
#endif

VM vm;

static Value
clock_native(uint8_t arg_count, Value *args)
{
  return NUMBER_VAL((double) clock() / CLOCKS_PER_SEC);
}

static void
vm_stack_reset()
{
  vm.stack_top = vm.stack;
  vm.frame_count = 0;
  vm.open_upvalues = NULL;
}

static void
runtime_error(const char *format, ...)
{
  va_list args;
  va_start(args, format);
  vfprintf(stderr, format, args);
  va_end(args);
  fputs("\n", stderr);
  for (int i = vm.frame_count - 1; i >= 0; --i) {
    CallFrame *frame = &vm.frames[i];
    ObjFunction *function = frame->closure->function;
    size_t instruction = frame->ip - function->chunk.code - 1;
    fprintf(stderr, "[line %d] in ", function->chunk.lines[instruction]);
    if (function->name == NULL)
      fprintf(stderr, "script\n");
    else
      fprintf(stderr, "%s()\n", function->name->chars);
  }
  vm_stack_reset();
}

static void
define_native(const char *name, NativeFn function)
{
  vm_stack_push(OBJ_VAL(copy_string(name, strlen(name))));
  vm_stack_push(OBJ_VAL(new_native(function)));
  table_set(&vm.globals, AS_STRING(vm.stack[0]), vm.stack[1]);
  vm_stack_pop();
  vm_stack_pop();
}

void
vm_init()
{
  vm_stack_reset();
  vm.objects = NULL;
  vm.bytes_allocated = 0;
  vm.next_gc = 1024 * 1024;
  vm.gray_count = 0;
  vm.gray_capacity = 0;
  vm.gray_stack = NULL;
  table_init(&vm.globals);
  table_init(&vm.strings);
  define_native("clock", clock_native);
}

void
vm_free()
{
  free_objects();
  table_free(&vm.globals);
  table_free(&vm.strings);
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
call(ObjClosure *closure, uint8_t arg_count)
{
  if (arg_count != closure->function->arity) {
    runtime_error("Expected %u arguments but got %u.",
        closure->function->arity, arg_count);
    return false;
  }
  if (vm.frame_count == FRAMES_MAX) {
    runtime_error("Stack overflow.");
    return false;
  }
  CallFrame *frame = &vm.frames[vm.frame_count++];
  frame->closure = closure;
  frame->ip = closure->function->chunk.code;
  frame->slots = vm.stack_top - arg_count - 1;
  return true;
}

static bool
call_value(Value callee, uint8_t arg_count)
{
  if (IS_OBJ(callee)) {
    switch (OBJ_TYPE(callee)) {
    case OBJ_CLASS: {
      ObjClass *class = AS_CLASS(callee);
      vm.stack_top[-arg_count - 1] = OBJ_VAL(new_instance(class));
      return true;
    }
    case OBJ_CLOSURE:
      return call(AS_CLOSURE(callee), arg_count);
    case OBJ_NATIVE: {
      NativeFn native = AS_NATIVE(callee);
      Value result = native(arg_count, vm.stack_top - arg_count);
      vm.stack_top -= arg_count + 1;
      vm_stack_push(result);
      return true;
    }
    default:
      break;
    }
  }
  runtime_error("Can only call functions and classes.");
  return false;
}

static ObjUpvalue *
capture_upvalue(Value *local)
{
  ObjUpvalue *prev_upvalue = NULL;
  ObjUpvalue *upvalue = vm.open_upvalues;
  while (upvalue != NULL && upvalue->location > local) {
    prev_upvalue = upvalue;
    upvalue = upvalue->next;
  }
  if (upvalue != NULL && upvalue->location == local)
    return upvalue;
  ObjUpvalue *created_upvalue = new_upvalue(local);
  created_upvalue->next = upvalue;
  if (prev_upvalue == NULL)
    vm.open_upvalues = created_upvalue;
  else
    prev_upvalue->next = created_upvalue;
  return created_upvalue;
}

static void
close_upvalues(Value *last)
{
  while (vm.open_upvalues != NULL && vm.open_upvalues->location >= last) {
    ObjUpvalue *upvalue = vm.open_upvalues;
    upvalue->closed = *upvalue->location;
    upvalue->location = &upvalue->closed;
    vm.open_upvalues = upvalue->next;
  }
}

static bool
is_falsey(Value value) {
  return IS_NIL(value) || (IS_BOOL(value) && !AS_BOOL(value));
}

static void
concatenate()
{
  ObjString *b = AS_STRING(vm_stack_peek(0));
  ObjString *a = AS_STRING(vm_stack_peek(1));
  size_t length = a->length + b->length;
  char *chars = ALLOCATE(char, length + 1);
  memcpy(chars, a->chars, a->length);
  memcpy(chars + a->length, b->chars, b->length);
  chars[length] = '\0';
  ObjString *result = take_string(chars, length);
  vm_stack_pop();
  vm_stack_pop();
  vm_stack_push(OBJ_VAL(result));
}

static InterpretResult
run()
{
  CallFrame *frame = &vm.frames[vm.frame_count - 1];
  #define READ_BYTE() (*frame->ip++)
  #define READ_SHORT() (frame->ip += 2, (uint16_t) ((frame->ip[-2] << 8) | frame->ip[-1]))
  #define READ_CONSTANT() (frame->closure->function->chunk.constants.values[READ_BYTE()])
  #define READ_STRING() AS_STRING(READ_CONSTANT())
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
    disassemble_instruction(&frame->closure->function->chunk,
        (size_t) (frame->ip - frame->closure->function->chunk.code));
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
    case OP_POP:
      vm_stack_pop();
      break;
    case OP_GET_LOCAL: {
      uint8_t slot = READ_BYTE();
      vm_stack_push(frame->slots[slot]);
      break;
    }
    case OP_SET_LOCAL: {
      uint8_t slot = READ_BYTE();
      frame->slots[slot] = vm_stack_peek(0);
      break;
    }
    case OP_GET_GLOBAL: {
      ObjString *name = READ_STRING();
      Value value;
      if (!table_get(&vm.globals, name, &value)) {
        runtime_error("Undefined variable '%s'.", name->chars);
        return INTERPRET_RUNTIME_ERROR;
      }
      vm_stack_push(value);
      break;
    }
    case OP_DEFINE_GLOBAL: {
      ObjString *name = READ_STRING();
      table_set(&vm.globals, name, vm_stack_peek(0));
      vm_stack_pop();
      break;
    }
    case OP_SET_GLOBAL: {
      ObjString *name = READ_STRING();
      if (table_set(&vm.globals, name, vm_stack_peek(0))) {
        table_delete(&vm.globals, name);
        runtime_error("Undefined variable '%s'.", name->chars);
        return INTERPRET_RUNTIME_ERROR;
      }
      break;
    }
    case OP_GET_UPVALUE: {
      uint8_t slot = READ_BYTE();
      vm_stack_push(*frame->closure->upvalues[slot]->location);
      break;
    }
    case OP_SET_UPVALUE: {
      uint8_t slot = READ_BYTE();
      *frame->closure->upvalues[slot]->location = vm_stack_peek(0);
      break;
    }
    case OP_GET_PROPERTY: {
      if (!IS_INSTANCE(vm_stack_peek(0))) {
        runtime_error("Only instances have properties.");
        return INTERPRET_RUNTIME_ERROR;
      }
      ObjInstance *instance = AS_INSTANCE(vm_stack_peek(0));
      ObjString *name = READ_STRING();
      Value value;
      if (table_get(&instance->fields, name, &value)) {
        vm_stack_pop();
        vm_stack_push(value);
        break;
      }
      runtime_error("Undefined property '%s'.", name->chars);
      return INTERPRET_RUNTIME_ERROR;
    }
    case OP_SET_PROPERTY: {
      if (!IS_INSTANCE(vm_stack_peek(1))) {
        runtime_error("Only instances have fields.");
        return INTERPRET_RUNTIME_ERROR;
      }
      ObjInstance *instance = AS_INSTANCE(vm_stack_peek(1));
      table_set(&instance->fields, READ_STRING(), vm_stack_peek(0));
      Value value = vm_stack_pop();
      vm_stack_pop();
      vm_stack_push(value);
      break;
    }
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
    case OP_PRINT:
      value_print(vm_stack_pop());
      printf("\n");
      break;
    case OP_JUMP: {
      uint16_t offset = READ_SHORT();
      frame->ip += offset;
      break;
    }
    case OP_JUMP_IF_FALSE: {
      uint16_t offset = READ_SHORT();
      if (is_falsey(vm_stack_peek(0)))
        frame->ip += offset;
      break;
    }
    case OP_LOOP: {
      uint16_t offset = READ_SHORT();
      frame->ip -= offset;
      break;
    }
    case OP_CALL: {
      uint8_t arg_count = READ_BYTE();
      if (!call_value(vm_stack_peek(arg_count), arg_count))
        return INTERPRET_RUNTIME_ERROR;
      frame = &vm.frames[vm.frame_count - 1];
      break;
    }
    case OP_CLOSURE: {
      ObjFunction *function = AS_FUNCTION(READ_CONSTANT());
      ObjClosure *closure = new_closure(function);
      vm_stack_push(OBJ_VAL(closure));
      for (size_t i = 0; i < closure->upvalue_count; ++i) {
        uint8_t is_local = READ_BYTE();
        uint8_t index = READ_BYTE();
        if (is_local)
          closure->upvalues[i] = capture_upvalue(frame->slots + index);
        else
          closure->upvalues[i] = frame->closure->upvalues[index];
      }
      break;
    }
    case OP_CLOSE_UPVALUE:
      close_upvalues(vm.stack_top - 1);
      vm_stack_pop();
      break;
    case OP_RETURN: {
      Value result = vm_stack_pop();
      close_upvalues(frame->slots);
      vm.frame_count--;
      if (vm.frame_count == 0) {
        vm_stack_pop();
        return INTERPRET_OK;
      }
      vm.stack_top = frame->slots;
      vm_stack_push(result);
      frame = &vm.frames[vm.frame_count - 1];
      break;
    }
    case OP_CLASS:
      vm_stack_push(OBJ_VAL(new_class(READ_STRING())));
      break;
    }
  }
  #undef READ_BYTE
  #undef READ_SHORT
  #undef READ_CONSTANT
  #undef READ_STRING
  #undef BINARY_OP
}

InterpretResult
vm_interpret(const char *source)
{
  ObjFunction *function = compile(source);
  if (function == NULL)
    return INTERPRET_COMPILE_ERROR;
  vm_stack_push(OBJ_VAL(function));
  ObjClosure *closure = new_closure(function);
  vm_stack_pop();
  vm_stack_push(OBJ_VAL(closure));
  call(closure, 0);
  return run();
}
