#ifndef CLOX_OBJECT_H
#define CLOX_OBJECT_H

#include "common.h"
#include "value.h"

#define OBJ_TYPE(value) (AS_OBJ(value)->type)
#define IS_STRING(value) is_obj_type(value, OBJ_STRING)

#define AS_STRING(value) ((ObjString *) AS_OBJ(value))
#define AS_CSTRING(value) (((ObjString *) AS_OBJ(value))->chars)

typedef enum {
  OBJ_STRING
} ObjType;

struct Obj {
  ObjType type;
  struct Obj *next;
};

struct ObjString {
  Obj obj;
  size_t length;
  char *chars;
};

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
