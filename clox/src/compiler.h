#ifndef CLOX_COMPILER_H
#define CLOX_COMPILER_H

#include "object.h"
#include "vm.h"

ObjFunction *
compile(const char *source);

void
compiler_mark_roots();

#endif
