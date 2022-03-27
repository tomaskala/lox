#ifndef CLOX_VALUE_H
#define CLOX_VALUE_H

#include "common.h"

typedef double Value;

typedef struct {
  size_t count;
  size_t capacity;
  Value *values;
} ValueArray;

void
value_array_init(ValueArray *array);

void
value_array_write(ValueArray *array, Value value);

void
value_array_free(ValueArray *array);

void
value_print(Value value);

#endif
