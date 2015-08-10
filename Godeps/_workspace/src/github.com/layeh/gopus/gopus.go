package gopus

// #include <opus.h>
//
// enum {
//   gopus_ok = OPUS_OK,
//   gopus_bad_arg = OPUS_BAD_ARG,
//   gopus_small_buffer = OPUS_BUFFER_TOO_SMALL,
//   gopus_internal = OPUS_INTERNAL_ERROR,
//   gopus_invalid_packet = OPUS_INVALID_PACKET,
//   gopus_unimplemented = OPUS_UNIMPLEMENTED,
//   gopus_invalid_state = OPUS_INVALID_STATE,
//   gopus_alloc_fail = OPUS_ALLOC_FAIL,
// };
import "C"

import (
	"errors"
)

var (
	ErrBadArgument   = errors.New("bad argument")
	ErrSmallBuffer   = errors.New("buffer is too small")
	ErrInternal      = errors.New("internal error")
	ErrInvalidPacket = errors.New("invalid packet")
	ErrUnimplemented = errors.New("unimplemented")
	ErrInvalidState  = errors.New("invalid state")
	ErrAllocFail     = errors.New("allocation failed")
	ErrUnknown       = errors.New("unknown error")
)

func getErr(code C.int) error {
	switch code {
	case C.gopus_ok:
		return nil
	case C.gopus_bad_arg:
		return ErrBadArgument
	case C.gopus_small_buffer:
		return ErrSmallBuffer
	case C.gopus_internal:
		return ErrInternal
	case C.gopus_invalid_packet:
		return ErrInvalidPacket
	case C.gopus_unimplemented:
		return ErrUnimplemented
	case C.gopus_invalid_state:
		return ErrInvalidState
	case C.gopus_alloc_fail:
		return ErrAllocFail
	default:
		return ErrUnknown
	}
}
