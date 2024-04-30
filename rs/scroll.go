package rs

/*
#cgo LDFLAGS: -L. -lm -ldl -lscroll_zstd
#include <stdlib.h>
#include "./rs_zstd.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
) //nolint:typecheck

func CompressScrollBatchBytes(src []byte) ([]byte, error) {

	srcSize := C.uint64_t(len(src))

	outbuf := make([]byte, len(src))

	outbufSize := C.uint64_t(len(src))

	err := C.compress_scroll_batch_bytes(
		(*C.uchar)(unsafe.Pointer(&src[0])),
		srcSize,
		(*C.uchar)(unsafe.Pointer(&outbuf[0])),
		&outbufSize,
	)

	if err != nil {
		return nil, fmt.Errorf("compress fail: %s", C.GoString(err))
	}

	return outbuf[:int(outbufSize)], nil
}
