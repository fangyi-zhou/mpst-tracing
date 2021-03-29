package pedro

// #include "ocaml_binding.h"
// #cgo LDFLAGS: -ldl
import "C"
import (
	"errors"
	"unsafe"
)

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

func (*OcamlRuntime) LoadFromFile(filename string) error {
	ret := C.pedro_load_from_file(C.CString(filename))
	if ret == nil {
		return nil
	} else {
		errMsg := C.GoString(ret)
		C.free(unsafe.Pointer(ret))
		return errors.New(errMsg)
	}
}

func (*OcamlRuntime) ImportNuscrFile(filename string, protocolName string) error {
	ret := C.pedro_import_nuscr_file(C.CString(filename), C.CString(protocolName))
	if ret == nil {
		return nil
	} else {
		errMsg := C.GoString(ret)
		C.free(unsafe.Pointer(ret))
		return errors.New(errMsg)
	}
}

func (*OcamlRuntime) SaveToFile(filename string) error {
	ret := C.pedro_save_to_file(C.CString(filename))
	if ret != 0 {
		return nil
	} else {
		return errors.New("unable to save to file")
	}
}

func (*OcamlRuntime) DoTransition(transition string) error {
	ret := C.pedro_do_transition(C.CString(transition))
	if ret != 0 {
		return nil
	} else {
		return errors.New("unable to do transition")
	}
}

func (*OcamlRuntime) GetEnabledTransitions() []string {
	var ret C.string_array_t
	C.pedro_get_enabled_transitions(&ret)
	retSize := ret.size
	retArray := ret.data
	i := 0
	size := int(retSize)
	transitions := make([]string, size)
	ptr := unsafe.Pointer(retArray)
	for i < size {
		transitions[i] = C.GoString(*(**C.char)(ptr))
		C.free(unsafe.Pointer(*(**C.char)(ptr)))
		ptr = unsafe.Pointer(uintptr(ptr) + unsafe.Sizeof(ptr))
		i++
	}
	C.free(unsafe.Pointer(retArray))
	return transitions
}

func (*OcamlRuntime) HasFinished() bool {
	return C.pedro_has_finished() == 1
}
