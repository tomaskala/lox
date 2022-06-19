#include <stdio.h>

#include "memory.h"
#include "value.h"

void
value_array_init(ValueArray *array)
{
  array->count = 0;
  array->capacity = 0;
  array->values = NULL;
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
  }
}
