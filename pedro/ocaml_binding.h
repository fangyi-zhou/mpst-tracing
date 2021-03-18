#ifndef PEDRO_BINDING_H
#define PEDRO_BINDING_H

#include <stdint.h>

// FIXME: This does not work across all archs
typedef intptr_t value;

typedef struct {
  void *handle; // shared object handle;
  void (*caml_startup)(char **argv);
  void (*caml_shutdown)(void);
  value (*caml_callback)(value, value);
  value *(*caml_named_value)(char const *);
  value (*caml_copy_string)(char const *);
} pedro_binding_t;

char *pedro_binding_init(char *);
void pedro_binding_deinit(void);
void pedro_call_main(char *);

#endif // PEDRO_BINDING_H
