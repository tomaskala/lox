#ifndef CLOX_VALUE_H
#define CLOX_VALUE_H

#include "common.h"

typedef struct Obj Obj;
typedef struct ObjString ObjString;

#ifdef NAN_BOXING

#include <string.h>

typedef uint64_t Value;

#define SIGN_BIT ((uint64_t) 0x8000000000000000)
#define QNAN ((uint64_t) 0x7ffc000000000000)
#define TAG_NIL 1
#define TAG_FALSE 2
#define TAG_TRUE 3

#define IS_BOOL(value) (((value) | 1) == TRUE_VAL)
#define IS_NIL(value) ((value) == NIL_VAL)
#define IS_NUMBER(value) (((value) & QNAN) != QNAN)
#define IS_OBJ(value) (((value) & (QNAN | SIGN_BIT)) == (QNAN | SIGN_BIT))

#define AS_BOOL(value) ((value) == TRUE_VAL)
#define AS_NUMBER(value) value_to_num(value)
#define AS_OBJ(value) ((Obj *) (uintptr_t) ((value) & ~(SIGN_BIT | QNAN)))

#define BOOL_VAL(b) ((b) ? TRUE_VAL : FALSE_VAL)
#define FALSE_VAL ((Value) (uint64_t) (QNAN | TAG_FALSE))
#define TRUE_VAL ((Value) (uint64_t) (QNAN | TAG_TRUE))
#define NIL_VAL ((Value) (uint64_t) (QNAN | TAG_NIL))
#define NUMBER_VAL(n) num_to_value(n)
#define OBJ_VAL(o) (Value) (SIGN_BIT | QNAN | (uint64_t) (uintptr_t) (o))

static inline double
value_to_num(Value value)
{
  double num;
  memcpy(&num, &value, sizeof(Value));
  return num;
}

static inline Value
num_to_value(double num)
{
  Value value;
  memcpy(&value, &num, sizeof(double));
  return value;
}

#else

typedef enum {
  VAL_BOOL,
  VAL_NIL,
  VAL_NUMBER,
  VAL_OBJ,
} ValueType;

typedef struct {
  ValueType type;
  union {
    bool boolean;
    double number;
    Obj *obj;
  } as;
} Value;

#define IS_BOOL(value) ((value).type == VAL_BOOL)
#define IS_NIL(value) ((value).type == VAL_NIL)
#define IS_NUMBER(value) ((value).type == VAL_NUMBER)
#define IS_OBJ(value) ((value).type == VAL_OBJ)

#define AS_BOOL(value) ((value).as.boolean)
#define AS_NUMBER(value) ((value).as.number)
#define AS_OBJ(value) ((value).as.obj)

#define BOOL_VAL(b) ((Value) {VAL_BOOL, {.boolean = b}})
#define NIL_VAL ((Value) {VAL_NIL, {.number = 0}})
#define NUMBER_VAL(n) ((Value) {VAL_NUMBER, {.number = n}})
#define OBJ_VAL(o) ((Value) {VAL_OBJ, {.obj = (Obj *) o}})

#endif

typedef struct {
  size_t count;
  size_t capacity;
  Value *values;
} ValueArray;

void
value_array_init(ValueArray *array);

bool values_equal(Value a, Value b);

void
value_array_write(ValueArray *array, Value value);

void
value_array_free(ValueArray *array);

void
value_print(Value value);

#endif
