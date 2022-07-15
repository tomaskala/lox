#ifndef CLOX_OBJECT_H
#define CLOX_OBJECT_H

#include "chunk.h"
#include "common.h"
#include "value.h"

#define OBJ_TYPE(value) (AS_OBJ(value)->type)
#define IS_FUNCTION(value) is_obj_type(value, OBJ_FUNCTION)
#define IS_STRING(value) is_obj_type(value, OBJ_STRING)

#define AS_FUNCTION(value) ((ObjFunction *) AS_OBJ(value))
#define AS_STRING(value) ((ObjString *) AS_OBJ(value))
#define AS_CSTRING(value) (((ObjString *) AS_OBJ(value))->chars)

typedef enum {
  OBJ_FUNCTION,
  OBJ_STRING,
} ObjType;

struct Obj {
  ObjType type;
  struct Obj *next;
};

typedef struct {
  Obj obj;
  size_t arity;
  Chunk chunk;
  ObjString *name;
} ObjFunction;

struct ObjString {
  Obj obj;
  size_t length;
  char *chars;
  uint32_t hash;
};

ObjFunction *
new_function();

ObjString *
take_string(char *chars, size_t length);

ObjString *
copy_string(const char *chars, size_t length);

void
object_print(Value value);

static inline bool
is_obj_type(Value value, ObjType type)
{
  return IS_OBJ(value) && OBJ_TYPE(value) == type;
}

#endif
