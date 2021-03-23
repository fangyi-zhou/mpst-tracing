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
  // shared object handle, loaded from `dlopen`
  void *handle;

  // Function pointers to OCaml runtime, loaded from `dlsym`
  FUNC_PTR(caml_startup);
  FUNC_PTR(caml_shutdown);
  FUNC_PTR(caml_callback);
  FUNC_PTR(caml_named_value);
  FUNC_PTR(caml_copy_string);

  // Pointers to exported OCaml functions
  // val main : string -> unit
  value main;
  // val load_from_file : string -> bool
  value load_from_file;
  // val save_to_file : string -> bool
  value save_to_file;
  // val get_enabled_transitions : unit -> string list
  value get_enabled_transitions;
  // val do_transition : string -> bool
  value do_transition;
} pedro_binding_t;

#undef FUNC_PTR

char *pedro_binding_init(char *);
void pedro_binding_deinit(void);
void pedro_call_main(char *);

#endif // PEDRO_BINDING_H
