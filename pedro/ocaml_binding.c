#include <dlfcn.h>
#include <stdlib.h>
#include <string.h>

#include "ocaml_binding.h"

static pedro_binding_t binding;

#define LOAD_SYM(SYMBOL)                                                       \
  SYMBOL##_t SYMBOL = dlsym(handle, #SYMBOL);                                  \
  if (!SYMBOL) {                                                               \
    return dlerror();                                                          \
  }                                                                            \
  binding.SYMBOL = SYMBOL;

#define LOAD_OCAML_VALUE(VALUE)                                                \
  value *VALUE = caml_named_value(#VALUE);                                     \
  if (!VALUE) {                                                                \
    return "Unable to get OCaml value " #VALUE;                                \
  }                                                                            \
  binding.VALUE = *VALUE;

// Returns an error string in case of failure, NULL in case of success
char *pedro_binding_init(char *path) {
  // Load shared object handle
  void *handle = dlopen(path, RTLD_LAZY);
  if (!handle) {
    return dlerror();
  }
  binding.handle = handle;

  // Load function symbols
  LOAD_SYM(caml_startup);
  LOAD_SYM(caml_shutdown);
  LOAD_SYM(caml_callback);
  LOAD_SYM(caml_named_value);
  LOAD_SYM(caml_copy_string);

  // Initialise OCaml runtime
  char *argv[1] = {NULL};
  caml_startup(argv);

  LOAD_OCAML_VALUE(main);
  LOAD_OCAML_VALUE(load_from_file);
  LOAD_OCAML_VALUE(save_to_file);
  LOAD_OCAML_VALUE(get_enabled_transitions);
  LOAD_OCAML_VALUE(do_transition);

  return NULL;
}

void pedro_binding_deinit(void) {
  if (!binding.handle) {
    return;
  }
  binding.caml_shutdown();
  memset(&binding, 0, sizeof(pedro_binding_t));
}

void pedro_call_main(char *filename) {
  static const value *main_closure = NULL;
  if (main_closure == NULL) {
    main_closure = binding.caml_named_value("main");
  }
  binding.caml_callback(*main_closure, binding.caml_copy_string(filename));
}

#undef LOAD_SYM
#undef LOAD_OCAML_VALUE