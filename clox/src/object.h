#ifndef CLOX_OBJECT_H
#define CLOX_OBJECT_H

#include "chunk.h"
#include "common.h"
#include "value.h"

#define OBJ_TYPE(value) (AS_OBJ(value)->type)
#define IS_CLOSURE(value) is_obj_type(value, OBJ_CLOSURE)
#define IS_FUNCTION(value) is_obj_type(value, OBJ_FUNCTION)
#define IS_NATIVE(value) is_obj_type(value, OBJ_NATIVE)
#define IS_STRING(value) is_obj_type(value, OBJ_STRING)

#define AS_CLOSURE(value) ((ObjClosure *) AS_OBJ(value))
#define AS_FUNCTION(value) ((ObjFunction *) AS_OBJ(value))
#define AS_NATIVE(value) (((ObjNative *) AS_OBJ(value))->function)
#define AS_STRING(value) ((ObjString *) AS_OBJ(value))
#define AS_CSTRING(value) (((ObjString *) AS_OBJ(value))->chars)

typedef enum {
  OBJ_CLOSURE,
  OBJ_FUNCTION,
  OBJ_NATIVE,
  OBJ_STRING,
  OBJ_UPVALUE,
} ObjType;

struct Obj {
  ObjType type;
  struct Obj *next;
};

typedef struct {
  Obj obj;
  size_t arity;
  size_t upvalue_count;
  Chunk chunk;
  ObjString *name;
} ObjFunction;

typedef Value (*NativeFn) (uint8_t arg_count, Value *args);

typedef struct {
  Obj obj;
  NativeFn function;
} ObjNative;

struct ObjString {
  Obj obj;
  size_t length;
  char *chars;
  uint32_t hash;
};

typedef struct ObjUpvalue {
  Obj obj;
  Value *location;
  Value closed;
  struct ObjUpvalue *next;
} ObjUpvalue;

typedef struct {
  Obj obj;
  ObjFunction *function;
  ObjUpvalue **upvalues;
  size_t upvalue_count;
} ObjClosure;

ObjClosure *
new_closure(ObjFunction *function);

ObjFunction *
new_function();

ObjNative *
new_native(NativeFn function);

ObjString *
take_string(char *chars, size_t length);

ObjString *
copy_string(const char *chars, size_t length);

ObjUpvalue *
new_upvalue(Value *slot);

void
object_print(Value value);

static inline bool
is_obj_type(Value value, ObjType type)
{
  return IS_OBJ(value) && OBJ_TYPE(value) == type;
}

#endif
