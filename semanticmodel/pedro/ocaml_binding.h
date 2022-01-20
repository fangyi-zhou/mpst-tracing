#ifndef PEDRO_BINDING_H
#define PEDRO_BINDING_H

#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

#include "version.h"

// FIXME: This does not work across all archs
typedef intptr_t value;
typedef struct {
  char **data;
  size_t size;
} string_array_t;

typedef void (*caml_startup_t)(char **);
typedef void (*caml_shutdown_t)(void);
typedef value (*caml_callback_t)(value, value);
typedef value (*caml_callback2_t)(value, value, value);
typedef value *(*caml_named_value_t)(char const *);
typedef value (*caml_copy_string_t)(char const *);

#define FUNC_PTR(NAME) NAME##_t NAME;

typedef struct {
  // shared object handle, loaded from `dlopen`
  void *handle;

  // indicate whether OCaml has been initialised;
  bool caml_initialised;

  // Function pointers to OCaml runtime, loaded from `dlsym`
  FUNC_PTR(caml_startup);
  FUNC_PTR(caml_shutdown);
  FUNC_PTR(caml_callback);
  FUNC_PTR(caml_callback2);
  FUNC_PTR(caml_named_value);
  FUNC_PTR(caml_copy_string);

  // Pointers to exported OCaml functions
  // val import_nuscr_file : string -> string -> string option
  value import_nuscr_file;
  // val load_from_file : string -> string option
  value load_from_file;
  // val save_to_file : string -> bool
  value save_to_file;
  // val get_enabled_transitions : unit -> string list
  value get_enabled_transitions;
  // val do_transition : string -> bool
  value do_transition;
  // val has_finished : unit -> bool
  value has_finished;
  // val commit_hash : string
  value commit_hash;
  // val get_all_roles : unit -> string list
  value get_all_roles;
} pedro_binding_t;

#undef FUNC_PTR

char *pedro_binding_init(char *);
void pedro_binding_deinit(void);
char *pedro_import_nuscr_file(char *, char *);
char *pedro_load_from_file(char *);
int pedro_save_to_file(char *);
int pedro_do_transition(char *);
void pedro_get_enabled_transitions(string_array_t *t);
void pedro_get_all_roles(string_array_t *t);
int pedro_has_finished(void);

#endif // PEDRO_BINDING_H
