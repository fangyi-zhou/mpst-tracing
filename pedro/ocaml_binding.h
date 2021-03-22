#ifndef PEDRO_BINDING_H
#define PEDRO_BINDING_H

#include <stdint.h>

// FIXME: This does not work across all archs
typedef intptr_t value;

typedef void (*caml_startup_t)(char **);
typedef void (*caml_shutdown_t)(void);
typedef value (*caml_callback_t)(value, value);
typedef value *(*caml_named_value_t)(char const *);
typedef value (*caml_copy_string_t)(char const *);

#define FUNC_PTR(NAME) NAME##_t NAME;

typedef struct {
  void *handle; // shared object handle;
  FUNC_PTR(caml_startup);
  FUNC_PTR(caml_shutdown);
  FUNC_PTR(caml_callback);
  FUNC_PTR(caml_named_value);
  FUNC_PTR(caml_copy_string);
} pedro_binding_t;

#undef FUNC_PTR

char *pedro_binding_init(char *);
void pedro_binding_deinit(void);
void pedro_call_main(char *);

#endif // PEDRO_BINDING_H
