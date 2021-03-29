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
  LOAD_SYM(caml_callback2);
  LOAD_SYM(caml_named_value);
  LOAD_SYM(caml_copy_string);

  // Initialise OCaml runtime
  char *argv[1] = {NULL};
  caml_startup(argv);

  LOAD_OCAML_VALUE(import_nuscr_file);
  LOAD_OCAML_VALUE(load_from_file);
  LOAD_OCAML_VALUE(save_to_file);
  LOAD_OCAML_VALUE(get_enabled_transitions);
  LOAD_OCAML_VALUE(do_transition);
  LOAD_OCAML_VALUE(has_finished);
  LOAD_OCAML_VALUE(commit_hash);

  char *hash = (char *)*commit_hash;
  if (strcmp(hash, PEDRO_API_HASH)) {
    return "Pedrolib version mismatch, expect commit hash " PEDRO_API_HASH;
  }

  return NULL;
}

void pedro_binding_deinit(void) {
  if (!binding.handle) {
    return;
  }
  binding.caml_shutdown();
  memset(&binding, 0, sizeof(pedro_binding_t));
}

char *pedro_load_from_file(char *filename) {
  value ret = binding.caml_callback(binding.load_from_file,
                                    binding.caml_copy_string(filename));
  free(filename);
  // interpret the return value as a string option
  if (ret == 1) {
    // None
    return NULL;
  }
  value *object = (value *)ret;
  char *err_string = (char *)(object[0]);
  char *dup = strdup(err_string);
  if (!dup) {
    return "pedro_load_from_file: Unable to get error message";
  }
  return dup;
}

char *pedro_import_nuscr_file(char *filename, char *protoname) {
  value ret = binding.caml_callback2(binding.import_nuscr_file,
                                     binding.caml_copy_string(filename),
                                     binding.caml_copy_string(protoname));
  free(filename);
  free(protoname);
  // interpret the return value as a string option
  if (ret == 1) {
    // None
    return NULL;
  }
  value *object = (value *)ret;
  char *err_string = (char *)(object[0]);
  char *dup = strdup(err_string);
  if (!dup) {
    return "pedro_import_nuscr_file: Unable to get error message";
  }
  return dup;
}

int pedro_save_to_file(char *filename) {
  value ret = binding.caml_callback(binding.save_to_file,
                                    binding.caml_copy_string(filename));
  // interpret the return value as a boolean
  free(filename);
  return (ret >> 1) ? 1 : 0;
}

int pedro_do_transition(char *transition) {
  value ret = binding.caml_callback(binding.do_transition,
                                    binding.caml_copy_string(transition));
  free(transition);
  // interpret the return value as a boolean
  return (ret >> 1) ? 1 : 0;
}

void pedro_get_enabled_transitions(string_array_t *out) {
  // FIXME: Check memory allocation failures and handle them gracefully
  size_t ptr_buf_size = 4;
  size_t idx = 0;
  char **ptr_out = malloc(ptr_buf_size * sizeof(char *));
  // 1 is an unit
  value ret = binding.caml_callback(binding.get_enabled_transitions, 1);
  value i = ret;

  // 1 is the empty list
  while (i != 1) {
    // A list is an object with two fields
    value *object = (value *)i;
    // first field is the list head
    const char *list_val = (const char *)object[0];
    if (ptr_buf_size == idx) {
      // Enlarge the array if full
      ptr_buf_size *= 2;
      ptr_out = realloc(ptr_out, ptr_buf_size * sizeof(char *));
    }
    ptr_out[idx] = strdup(list_val);
    idx++;
    // second field is the list tail
    i = object[1];
  }
  out->data = ptr_out;
  out->size = idx;
}

int pedro_has_finished(void) {
  value ret = binding.caml_callback(binding.do_transition, 1);
  // interpret the return value as a boolean
  return (ret >> 1) ? 1 : 0;
}

#undef LOAD_SYM
#undef LOAD_OCAML_VALUE