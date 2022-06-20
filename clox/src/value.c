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
  if (a.type != b.type)
    return false;
  switch (a.type) {
  case VAL_BOOL:
    return AS_BOOL(a) == AS_BOOL(b);
  case VAL_NIL:
    return true;
  case VAL_NUMBER:
    return AS_NUMBER(a) == AS_NUMBER(b);
  case VAL_OBJ: {
    ObjString *a_string = AS_STRING(a);
    ObjString *b_string = AS_STRING(b);
    return a_string->length == b_string->length
      && memcmp(a_string->chars, b_string->chars, a_string->length) == 0;
  }
  default:
    // Unreachable.
    return false;
  }
}

void
value_array_write(ValueArray *array, Value value)
{
  if (array->capacity < array->count + 1) {
    size_t old_capacity = array->capacity;
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
}
