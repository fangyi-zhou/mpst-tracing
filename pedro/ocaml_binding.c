#include <dlfcn.h>
#include <stdlib.h>
#include <string.h>

#include "ocaml_binding.h"

static pedro_binding_t binding;

// Returns an error string in case of failure, NULL in case of success
char *pedro_binding_init(char *path) {
  // Load shared object handle
  void *handle = dlopen(path, RTLD_LAZY);
  if (!handle) {
    return dlerror();
  }

  // Load function symbols
  void (*caml_startup)(char **argv) = dlsym(handle, "caml_startup");
  if (!caml_startup) {
    return dlerror();
  }
  void (*caml_shutdown)(void) = dlsym(handle, "caml_shutdown");
  if (!caml_shutdown) {
    return dlerror();
  }
  value (*caml_callback)(value, value) = dlsym(handle, "caml_callback");
  if (!caml_callback) {
    return dlerror();
  }
  value *(*caml_named_value)(char const *) = dlsym(handle, "caml_named_value");
  if (!caml_named_value) {
    return dlerror();
  }
  value (*caml_copy_string)(char const *) = dlsym(handle, "caml_copy_string");
  if (!caml_copy_string) {
    return dlerror();
  }

  // Store the pointers to the binding
  binding.handle = handle;
  binding.caml_startup = caml_startup;
  binding.caml_shutdown = caml_shutdown;
  binding.caml_callback = caml_callback;
  binding.caml_named_value = caml_named_value;
  binding.caml_copy_string = caml_copy_string;

  // Initialise OCaml runtime
  char *argv[1] = {NULL};
  caml_startup(argv);

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
