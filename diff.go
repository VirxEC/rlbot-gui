package main

/*
#cgo windows CFLAGS: -I./bob-diff
#cgo windows LDFLAGS: ./bob-diff/libbob_diff_win.a -lntdll -luserenv
#cgo linux LDFLAGS: ./bob-diff/libbob_diff.a -ldl
#include "./bob-diff/bob_diff.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

// func diff(old_path string, new_path string) []byte {
// 	old_path_c := C.CString(old_path)
// 	new_path_c := C.CString(new_path)
// 	defer C.free(unsafe.Pointer(old_path_c))
// 	defer C.free(unsafe.Pointer(new_path_c))

// 	buf := (*C.char)(nil)
// 	buf_len := C.uint32_t(0)

// 	err := C.diff(old_path_c, new_path_c, &buf, &buf_len)
// 	if err != 0 {
// 		return nil
// 	}

// 	return C.GoBytes(unsafe.Pointer(buf), C.int(buf_len))
// }

func diffApply(path string, diff []byte) error {
	pathC := C.CString(path)
	defer C.free(unsafe.Pointer(pathC))

	// buff is to the first byte of the slice
	buf := (*C.char)(unsafe.Pointer(&diff[0]))
	bufLen := C.uint32_t(len(diff))

	err := C.diff_apply(pathC, buf, bufLen)
	if err != 0 {
		return errors.New("failed to apply diff")
	}

	return nil
}
