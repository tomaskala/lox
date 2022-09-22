#ifndef CLOX_TABLE_H
#define CLOX_TABLE_H

#include "common.h"
#include "value.h"

typedef struct {
  ObjString *key;
  Value value;
} Entry;

typedef struct {
  int count;
  int capacity;
  Entry *entries;
} Table;

void
table_init(Table *table);

void
table_free(Table *table);

bool
table_get(Table *table, ObjString *key, Value *value);

bool
table_set(Table *table, ObjString *key, Value value);

bool
table_delete(Table *table, ObjString *key);

void
table_add_all(Table *from, Table *to);

ObjString *
table_find_string(Table *table, const char *chars, int length, uint32_t hash);

void
table_remove_white(Table *table);

void
table_mark(Table *table);

#endif
