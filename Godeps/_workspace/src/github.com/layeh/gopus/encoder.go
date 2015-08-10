package gopus

// #cgo !nopkgconfig pkg-config: opus
// #include <opus.h>
//
// enum {
//   gopus_application_voip    = OPUS_APPLICATION_VOIP,
//   gopus_application_audio   = OPUS_APPLICATION_AUDIO,
//   gopus_restricted_lowdelay = OPUS_APPLICATION_RESTRICTED_LOWDELAY,
//   gopus_bitrate_max         = OPUS_BITRATE_MAX,
// };
//
//
// void gopus_setvbr(OpusEncoder *encoder, int vbr) {
//   opus_encoder_ctl(encoder, OPUS_SET_VBR(vbr));
// }
//
// void gopus_setbitrate(OpusEncoder *encoder, int bitrate) {
//   opus_encoder_ctl(encoder, OPUS_SET_BITRATE(bitrate));
// }
//
// opus_int32 gopus_bitrate(OpusEncoder *encoder) {
//   opus_int32 bitrate;
//   opus_encoder_ctl(encoder, OPUS_GET_BITRATE(&bitrate));
//   return bitrate;
// }
//
// void gopus_setapplication(OpusEncoder *encoder, int application) {
//   opus_encoder_ctl(encoder, OPUS_SET_APPLICATION(application));
// }
//
// opus_int32 gopus_application(OpusEncoder *encoder) {
//   opus_int32 application;
//   opus_encoder_ctl(encoder, OPUS_GET_APPLICATION(&application));
//   return application;
// }
//
// void gopus_encoder_resetstate(OpusEncoder *encoder) {
//   opus_encoder_ctl(encoder, OPUS_RESET_STATE);
// }
import "C"

import (
	"unsafe"
)

type Application int

const (
	Voip               Application = C.gopus_application_voip
	Audio              Application = C.gopus_application_audio
	RestrictedLowDelay Application = C.gopus_restricted_lowdelay
)

const (
	BitrateMaximum = C.gopus_bitrate_max
)

type Encoder struct {
	data     []byte
	cEncoder *C.struct_OpusEncoder
}

func NewEncoder(sampleRate, channels int, application Application) (*Encoder, error) {
	encoder := &Encoder{}
	encoder.data = make([]byte, int(C.opus_encoder_get_size(C.int(channels))))
	encoder.cEncoder = (*C.struct_OpusEncoder)(unsafe.Pointer(&encoder.data[0]))

	ret := C.opus_encoder_init(encoder.cEncoder, C.opus_int32(sampleRate), C.int(channels), C.int(application))
	if err := getErr(ret); err != nil {
		return nil, err
	}
	return encoder, nil
}

func (e *Encoder) Encode(pcm []int16, frameSize, maxDataBytes int) ([]byte, error) {
	pcmPtr := (*C.opus_int16)(unsafe.Pointer(&pcm[0]))

	data := make([]byte, maxDataBytes)
	dataPtr := (*C.uchar)(unsafe.Pointer(&data[0]))

	encodedC := C.opus_encode(e.cEncoder, pcmPtr, C.int(frameSize), dataPtr, C.opus_int32(len(data)))
	encoded := int(encodedC)

	if encoded < 0 {
		return nil, getErr(C.int(encodedC))
	}
	return data[0:encoded], nil
}

func (e *Encoder) SetVbr(vbr bool) {
	var cVbr C.int
	if vbr {
		cVbr = 1
	} else {
		cVbr = 0
	}
	C.gopus_setvbr(e.cEncoder, cVbr)
}

func (e *Encoder) SetBitrate(bitrate int) {
	C.gopus_setbitrate(e.cEncoder, C.int(bitrate))
}

func (e *Encoder) Bitrate() int {
	return int(C.gopus_bitrate(e.cEncoder))
}

func (e *Encoder) SetApplication(application Application) {
	C.gopus_setapplication(e.cEncoder, C.int(application))
}

func (e *Encoder) Application() Application {
	return Application(C.gopus_application(e.cEncoder))
}

func (e *Encoder) ResetState() {
	C.gopus_encoder_resetstate(e.cEncoder)
}
