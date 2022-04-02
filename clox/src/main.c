#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "chunk.h"
#include "vm.h"

static void
repl()
{
  char line[1024];
  for (;;) {
    printf("> ");
    if (!fgets(line, sizeof(line), stdin)) {
      printf("\n");
      break;
    }
    vm_interpret(line);
  }
}

static char *
read_file(const char *path)
{
  FILE *file = fopen(path, "rb");
  if (file == NULL) {
    fprintf(stderr, "Could not open file \"%s\".\n", path);
    return NULL;
  }
  fseek(file, 0L, SEEK_END);
  size_t file_size = ftell(file);
  rewind(file);
  char *buffer = malloc(file_size + 1);
  if (buffer == NULL) {
    fprintf(stderr, "Not enough memory to read \"%s\".\n", path);
    return NULL;
  }
  size_t bytes_read = fread(buffer, sizeof(char), file_size, file);
  if (bytes_read < file_size) {
    free(buffer);
    fprintf(stderr, "Could not read file \"%s\".\n", path);
    return NULL;
  }
  buffer[bytes_read] = '\0';
  fclose(file);
  return buffer;
}

static int
run_file(const char *path)
{
  char *source = read_file(path);
  if (source == NULL)
    return 74;
  InterpretResult result = vm_interpret(source);
  free(source);
  if (result == INTERPRET_COMPILE_ERROR)
    return 65;
  else if (result == INTERPRET_RUNTIME_ERROR)
    return 70;
  else
    return 0;
}

int
main(int argc, const char **argv)
{
  int return_code = 0;
  vm_init();
  if (argc == 1)
    repl();
  else if (argc == 2)
    return_code = run_file(argv[1]);
  else {
    fprintf(stderr, "Usage: clox [path]\n");
    return_code = 64;
  }
  vm_free();
  return return_code;
}
