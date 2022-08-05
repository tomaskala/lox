#include <stdlib.h>

#include "compiler.h"
#include "memory.h"
#include "vm.h"

#ifdef DEBUG_LOG_GC
#include <stdio.h>
#include "debug.h"
#endif

#define GC_HEAP_GROW_FACTOR 2

void *
reallocate(void *pointer, size_t old_size, size_t new_size)
{
  vm.bytes_allocated += new_size - old_size;
  if (new_size > old_size) {
    #ifdef DEBUG_STRESS_GC
    collect_garbage();
    #endif
    if (vm.bytes_allocated > vm.next_gc)
      collect_garbage();
  }
  if (new_size == 0) {
    free(pointer);
    return NULL;
  }
  void *result = realloc(pointer, new_size);
  if (result == NULL)
    exit(1);
  return result;
}

void
object_mark(Obj *object)
{
  if (object == NULL)
    return;
  if (object->is_marked)
    return;
  #ifdef DEBUG_LOG_GC
  printf("%p mark ", (void *) object);
  value_print(OBJ_VAL(object));
  printf("\n");
  #endif
  object->is_marked = true;
  if (vm.gray_capacity < vm.gray_count + 1) {
    vm.gray_capacity = GROW_CAPACITY(vm.gray_capacity);
    vm.gray_stack = (Obj **) realloc(vm.gray_stack,
        sizeof(Obj *) * vm.gray_capacity);
    if (vm.gray_stack == NULL)
      exit(1);
  }
  vm.gray_stack[vm.gray_count++] = object;
}

void
value_mark(Value value)
{
  if (IS_OBJ(value))
    object_mark(AS_OBJ(value));
}

static void
array_mark(ValueArray *array)
{
  for (size_t i = 0; i < array->count; ++i)
    value_mark(array->values[i]);
}

static void
object_blacken(Obj *object)
{
  #ifdef DEBUG_LOG_GC
  printf("%p blacken ", (void *) object);
  value_print(OBJ_VAL(object));
  printf("\n");
  #endif
  switch (object->type) {
  case OBJ_BOUND_METHOD: {
    ObjBoundMethod *bound = (ObjBoundMethod *) object;
    value_mark(bound->receiver);
    object_mark((Obj *) bound->method);
    break;
  }
  case OBJ_CLASS: {
    ObjClass *class = (ObjClass *) object;
    object_mark((Obj *) class->name);
    table_mark(&class->methods);
    break;
  }
  case OBJ_CLOSURE: {
    ObjClosure *closure = (ObjClosure *) object;
    object_mark((Obj *) closure->function);
    for (size_t i = 0; i < closure->upvalue_count; ++i)
      object_mark((Obj *) closure->upvalues[i]);
    break;
  }
  case OBJ_FUNCTION: {
    ObjFunction *function = (ObjFunction *) object;
    object_mark((Obj *) function->name);
    array_mark(&function->chunk.constants);
    break;
  }
  case OBJ_INSTANCE: {
    ObjInstance *instance = (ObjInstance *) object;
    object_mark((Obj *) instance->class);
    table_mark(&instance->fields);
    break;
  }
  case OBJ_UPVALUE:
    value_mark(((ObjUpvalue *) object)->closed);
    break;
  case OBJ_NATIVE:
  case OBJ_STRING:
    break;
  }
}

static void
free_object(Obj *object)
{
  #ifdef DEBUG_LOG_GC
  printf("%p free type %d\n", (void *) object, object->type);
  #endif
  switch (object->type) {
  case OBJ_BOUND_METHOD:
    FREE(ObjBoundMethod, object);
    break;
  case OBJ_CLASS: {
    ObjClass *class = (ObjClass *) object;
    table_free(&class->methods);
    FREE(ObjClass, object);
    break;
  }
  case OBJ_CLOSURE: {
    ObjClosure *closure = (ObjClosure *) object;
    FREE_ARRAY(ObjUpvalue *, closure->upvalues, closure->upvalue_count);
    FREE(ObjClosure, object);
    break;
  }
  case OBJ_FUNCTION: {
    ObjFunction *function = (ObjFunction *) object;
    chunk_free(&function->chunk);
    FREE(ObjFunction, object);
    break;
  }
  case OBJ_INSTANCE: {
    ObjInstance *instance = (ObjInstance *) object;
    table_free(&instance->fields);
    FREE(ObjInstance, object);
    break;
  }
  case OBJ_NATIVE:
    FREE(ObjNative, object);
    break;
  case OBJ_STRING: {
    ObjString *string = (ObjString *) object;
    FREE_ARRAY(char, string->chars, string->length + 1);
    FREE(ObjString, object);
    break;
  }
  case OBJ_UPVALUE: {
    FREE(ObjUpvalue, object);
    break;
  }
  }
}

static void
mark_roots()
{
  for (Value *slot = vm.stack; slot < vm.stack_top; ++slot)
    value_mark(*slot);
  for (size_t i = 0; i < vm.frame_count; ++i)
    object_mark((Obj *) vm.frames[i].closure);
  for (ObjUpvalue *upvalue = vm.open_upvalues; upvalue != NULL;
      upvalue = upvalue->next)
    object_mark((Obj *) upvalue);
  table_mark(&vm.globals);
  compiler_mark_roots();
  object_mark((Obj *) vm.init_string);
}

static void
trace_references()
{
  while (vm.gray_count > 0) {
    Obj *object = vm.gray_stack[--vm.gray_count];
    object_blacken(object);
  }
}

static void
sweep()
{
  Obj *previous = NULL;
  Obj *object = vm.objects;
  while (object != NULL) {
    if (object->is_marked) {
      object->is_marked = false;
      previous = object;
      object = object->next;
    } else {
      Obj *unreached = object;
      object = object->next;
      if (previous != NULL)
        previous->next = object;
      else
        vm.objects = object;
      free_object(unreached);
    }
  }
}

void
collect_garbage()
{
  #ifdef DEBUG_LOG_GC
  printf("-- gc begin\n");
  size_t before = vm.bytes_allocated;
  #endif
  mark_roots();
  trace_references();
  table_remove_white(&vm.strings);
  sweep();
  vm.next_gc = vm.bytes_allocated * GC_HEAP_GROW_FACTOR;
  #ifdef DEBUG_LOG_GC
  printf("-- gc end\n");
  printf(" collected %zu bytes (from %zu to %zu) next at %zu\n",
      before - vm.bytes_allocated, before, vm.bytes_allocated, vm.next_gc);
  #endif
}

void
free_objects()
{
  Obj *object = vm.objects;
  while (object != NULL) {
    Obj *next = object->next;
    free_object(object);
    object = next;
  }
  free(vm.gray_stack);
}
