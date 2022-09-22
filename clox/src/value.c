#include <stdio.h>
#include <string.h>

#include "object.h"
#include "memory.h"
#include "value.h"

void
value_array_init(ValueArray *array)
{
  array->count = 0;
  array->capacity = 0;
  array->values = NULL;
}

bool
values_equal(Value a, Value b)
{
  #ifdef NAN_BOXING
  if (IS_NUMBER(a) && IS_NUMBER(b))
    return AS_NUMBER(a) == AS_NUMBER(b);
  return a == b;
  #else
  if (a.type != b.type)
    return false;
  switch (a.type) {
  case VAL_BOOL:
    return AS_BOOL(a) == AS_BOOL(b);
  case VAL_NIL:
    return true;
  case VAL_NUMBER:
    return AS_NUMBER(a) == AS_NUMBER(b);
  case VAL_OBJ:
    return AS_OBJ(a) == AS_OBJ(b);
  default:
    // Unreachable.
    return false;
  }
  #endif
}

void
value_array_write(ValueArray *array, Value value)
{
  if (array->capacity < array->count + 1) {
    int old_capacity = array->capacity;
    array->capacity = GROW_CAPACITY(old_capacity);
    array->values = GROW_ARRAY(Value, array->values, old_capacity,
                               array->capacity);
  }
  array->values[array->count++] = value;
}

void
value_array_free(ValueArray *array)
{
  FREE_ARRAY(Value, array->values, array->capacity);
  value_array_init(array);
}

void
value_print(Value value)
{
  #ifdef NAN_BOXING
  if (IS_BOOL(value))
    printf(AS_BOOL(value) ? "true" : "false");
  else if (IS_NIL(value))
    printf("nil");
  else if (IS_NUMBER(value))
    printf("%g", AS_NUMBER(value));
  else if (IS_OBJ(value))
    object_print(value);
  #else
  switch (value.type) {
  case VAL_BOOL:
    printf(AS_BOOL(value) ? "true" : "false");
    break;
  case VAL_NIL:
    printf("nil");
    break;
  case VAL_NUMBER:
    printf("%g", AS_NUMBER(value));
    break;
  case VAL_OBJ:
    object_print(value);
    break;
  }
  #endif
}
