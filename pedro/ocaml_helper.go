package pedro

// #include "ocaml_binding.h"
// #cgo LDFLAGS: -ldl
import "C"
import "errors"

type OcamlRuntime struct{}

func LoadRuntime(path string) (*OcamlRuntime, error) {
	errMsg := C.pedro_binding_init(C.CString(path))
	if errMsg != nil {
		errStr := C.GoString(errMsg)
		err := errors.New(errStr)
		return nil, err
	}
	return &OcamlRuntime{}, nil
}

func (*OcamlRuntime) Close() {
	C.pedro_binding_deinit()
}

func (*OcamlRuntime) RunMain(filename string) {
	C.pedro_call_main(C.CString(filename))
}

func (*OcamlRuntime) LoadFromFile(filename string) bool {
	ret := C.pedro_load_from_file(C.CString(filename))
	if ret == 0 {
		return false
	} else {
		return true
	}
}

func (*OcamlRuntime) SaveToFile(filename string) bool {
	ret := C.pedro_save_to_file(C.CString(filename))
	if ret == 0 {
		return false
	} else {
		return true
	}
}

func (*OcamlRuntime) DoTransition(transition string) bool {
	ret := C.pedro_do_transition(C.CString(transition))
	if ret == 0 {
		return false
	} else {
		return true
	}
}
