// +build amd64,cgo 386,cgo

package gopus // import "layeh.com/gopus"

// #cgo linux darwin freebsd LDFLAGS: -lm
//
// #cgo CFLAGS: -Iopus-1.1.2/include -Iopus-1.1.2/celt -Iopus-1.1.2/silk -Iopus-1.1.2/silk/float -Iopus-1.1.2/silk/fixed
//
// #include "opus-1.1.2/config.h"
//
// #include "opus-1.1.2/silk/CNG.c"
// #include "opus-1.1.2/silk/code_signs.c"
// #include "opus-1.1.2/silk/init_decoder.c"
// #include "opus-1.1.2/silk/decode_core.c"
// #include "opus-1.1.2/silk/decode_frame.c"
// #include "opus-1.1.2/silk/decode_parameters.c"
// #include "opus-1.1.2/silk/decode_indices.c"
// #include "opus-1.1.2/silk/decode_pulses.c"
// #include "opus-1.1.2/silk/decoder_set_fs.c"
// #include "opus-1.1.2/silk/dec_API.c"
// #include "opus-1.1.2/silk/enc_API.c"
// #include "opus-1.1.2/silk/encode_indices.c"
// #include "opus-1.1.2/silk/encode_pulses.c"
// #include "opus-1.1.2/silk/gain_quant.c"
// #include "opus-1.1.2/silk/interpolate.c"
// #include "opus-1.1.2/silk/LP_variable_cutoff.c"
// #include "opus-1.1.2/silk/NLSF_decode.c"
// #include "opus-1.1.2/silk/NSQ.c"
// #include "opus-1.1.2/silk/NSQ_del_dec.c"
// #include "opus-1.1.2/silk/PLC.c"
// #include "opus-1.1.2/silk/shell_coder.c"
// #include "opus-1.1.2/silk/tables_gain.c"
// #include "opus-1.1.2/silk/tables_LTP.c"
// #include "opus-1.1.2/silk/tables_NLSF_CB_NB_MB.c"
// #include "opus-1.1.2/silk/tables_NLSF_CB_WB.c"
// #include "opus-1.1.2/silk/tables_other.c"
// #include "opus-1.1.2/silk/tables_pitch_lag.c"
// #include "opus-1.1.2/silk/tables_pulses_per_block.c"
// #include "opus-1.1.2/silk/VAD.c"
// #include "opus-1.1.2/silk/control_audio_bandwidth.c"
// #include "opus-1.1.2/silk/quant_LTP_gains.c"
// #include "opus-1.1.2/silk/VQ_WMat_EC.c"
// #include "opus-1.1.2/silk/HP_variable_cutoff.c"
// #include "opus-1.1.2/silk/NLSF_encode.c"
// #include "opus-1.1.2/silk/NLSF_VQ.c"
// #include "opus-1.1.2/silk/NLSF_unpack.c"
// #include "opus-1.1.2/silk/NLSF_del_dec_quant.c"
// #include "opus-1.1.2/silk/process_NLSFs.c"
// #include "opus-1.1.2/silk/stereo_LR_to_MS.c"
// #include "opus-1.1.2/silk/stereo_MS_to_LR.c"
// #include "opus-1.1.2/silk/check_control_input.c"
// #include "opus-1.1.2/silk/control_SNR.c"
// #include "opus-1.1.2/silk/init_encoder.c"
// #include "opus-1.1.2/silk/control_codec.c"
// #include "opus-1.1.2/silk/A2NLSF.c"
// #include "opus-1.1.2/silk/ana_filt_bank_1.c"
// #include "opus-1.1.2/silk/biquad_alt.c"
// #include "opus-1.1.2/silk/bwexpander_32.c"
// #include "opus-1.1.2/silk/bwexpander.c"
// #include "opus-1.1.2/silk/debug.c"
// #include "opus-1.1.2/silk/decode_pitch.c"
// #include "opus-1.1.2/silk/inner_prod_aligned.c"
// #include "opus-1.1.2/silk/lin2log.c"
// #include "opus-1.1.2/silk/log2lin.c"
// #include "opus-1.1.2/silk/LPC_analysis_filter.c"
// #include "opus-1.1.2/silk/LPC_inv_pred_gain.c"
// #undef QA
// #include "opus-1.1.2/silk/table_LSF_cos.c"
// #include "opus-1.1.2/silk/NLSF2A.c"
// #undef QA
// #include "opus-1.1.2/silk/NLSF_stabilize.c"
// #include "opus-1.1.2/silk/NLSF_VQ_weights_laroia.c"
// #include "opus-1.1.2/silk/pitch_est_tables.c"
// #include "opus-1.1.2/silk/resampler.c"
// #include "opus-1.1.2/silk/resampler_down2_3.c"
// #include "opus-1.1.2/silk/resampler_down2.c"
// #include "opus-1.1.2/silk/resampler_private_AR2.c"
// #include "opus-1.1.2/silk/resampler_private_down_FIR.c"
// #include "opus-1.1.2/silk/resampler_private_IIR_FIR.c"
// #include "opus-1.1.2/silk/resampler_private_up2_HQ.c"
// #include "opus-1.1.2/silk/resampler_rom.c"
// #include "opus-1.1.2/silk/sigm_Q15.c"
// #include "opus-1.1.2/silk/sort.c"
// #include "opus-1.1.2/silk/sum_sqr_shift.c"
// #include "opus-1.1.2/silk/stereo_decode_pred.c"
// #include "opus-1.1.2/silk/stereo_encode_pred.c"
// #include "opus-1.1.2/silk/stereo_find_predictor.c"
// #include "opus-1.1.2/silk/stereo_quant_pred.c"
//
// #include "opus-1.1.2/silk/float/apply_sine_window_FLP.c"
// #include "opus-1.1.2/silk/float/corrMatrix_FLP.c"
// #include "opus-1.1.2/silk/float/encode_frame_FLP.c"
// #include "opus-1.1.2/silk/float/find_LPC_FLP.c"
// #include "opus-1.1.2/silk/float/find_LTP_FLP.c"
// #include "opus-1.1.2/silk/float/find_pitch_lags_FLP.c"
// #include "opus-1.1.2/silk/float/find_pred_coefs_FLP.c"
// #include "opus-1.1.2/silk/float/LPC_analysis_filter_FLP.c"
// #include "opus-1.1.2/silk/float/LTP_analysis_filter_FLP.c"
// #include "opus-1.1.2/silk/float/LTP_scale_ctrl_FLP.c"
// #include "opus-1.1.2/silk/float/noise_shape_analysis_FLP.c"
// #include "opus-1.1.2/silk/float/prefilter_FLP.c"
// #include "opus-1.1.2/silk/float/process_gains_FLP.c"
// #include "opus-1.1.2/silk/float/regularize_correlations_FLP.c"
// #include "opus-1.1.2/silk/float/residual_energy_FLP.c"
// #include "opus-1.1.2/silk/float/solve_LS_FLP.c"
// #include "opus-1.1.2/silk/float/warped_autocorrelation_FLP.c"
// #include "opus-1.1.2/silk/float/wrappers_FLP.c"
// #include "opus-1.1.2/silk/float/autocorrelation_FLP.c"
// #include "opus-1.1.2/silk/float/burg_modified_FLP.c"
// #include "opus-1.1.2/silk/float/bwexpander_FLP.c"
// #include "opus-1.1.2/silk/float/energy_FLP.c"
// #include "opus-1.1.2/silk/float/inner_product_FLP.c"
// #include "opus-1.1.2/silk/float/k2a_FLP.c"
// #include "opus-1.1.2/silk/float/levinsondurbin_FLP.c"
// #include "opus-1.1.2/silk/float/LPC_inv_pred_gain_FLP.c"
// #include "opus-1.1.2/silk/float/pitch_analysis_core_FLP.c"
// #include "opus-1.1.2/silk/float/scale_copy_vector_FLP.c"
// #include "opus-1.1.2/silk/float/scale_vector_FLP.c"
// #include "opus-1.1.2/silk/float/schur_FLP.c"
// #include "opus-1.1.2/silk/float/sort_FLP.c"
//
// #undef PI
// #include "opus-1.1.2/celt/bands.c"
// OPUS_CUSTOM_NOSTATIC int opus_custom_encoder_get_size(const CELTMode *mode, int channels);
// #include "opus-1.1.2/celt/celt.c"
// int opus_custom_encoder_init(CELTEncoder *st, const CELTMode *mode, int channels);
// #include "opus-1.1.2/celt/celt_encoder.c"
// OPUS_CUSTOM_NOSTATIC int opus_custom_decoder_get_size(const CELTMode *mode, int channels);
// OPUS_CUSTOM_NOSTATIC int opus_custom_decoder_init(CELTDecoder *st, const CELTMode *mode, int channels);
// #include "opus-1.1.2/celt/celt_decoder.c"
// #include "opus-1.1.2/celt/cwrs.c"
// #include "opus-1.1.2/celt/entcode.c"
// #include "opus-1.1.2/celt/entdec.c"
// #include "opus-1.1.2/celt/entenc.c"
// #include "opus-1.1.2/celt/kiss_fft.c"
// #include "opus-1.1.2/celt/laplace.c"
// #include "opus-1.1.2/celt/mathops.c"
// #include "opus-1.1.2/celt/mdct.c"
// #include "opus-1.1.2/celt/modes.c"
// #include "opus-1.1.2/celt/pitch.c"
// #include "opus-1.1.2/celt/celt_lpc.c"
// #include "opus-1.1.2/celt/quant_bands.c"
// #include "opus-1.1.2/celt/rate.c"
// #include "opus-1.1.2/celt/vq.c"
//
// #include "opus-1.1.2/src/opus.c"
// #include "opus-1.1.2/src/opus_decoder.c"
// #include "opus-1.1.2/src/opus_encoder.c"
// #include "opus-1.1.2/src/opus_multistream.c"
// #include "opus-1.1.2/src/opus_multistream_encoder.c"
// #include "opus-1.1.2/src/opus_multistream_decoder.c"
// #include "opus-1.1.2/src/repacketizer.c"
//
// #include "opus-1.1.2/src/analysis.c"
// #include "opus-1.1.2/src/mlp.c"
// #include "opus-1.1.2/src/mlp_data.c"
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
//
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
//
// void gopus_decoder_resetstate(OpusDecoder *decoder) {
//   opus_decoder_ctl(decoder, OPUS_RESET_STATE);
// }
import "C"

import (
	"errors"
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

type Decoder struct {
	data     []byte
	cDecoder *C.struct_OpusDecoder
	channels int
}

func NewDecoder(sampleRate, channels int) (*Decoder, error) {
	decoder := &Decoder{}
	decoder.data = make([]byte, int(C.opus_decoder_get_size(C.int(channels))))
	decoder.cDecoder = (*C.struct_OpusDecoder)(unsafe.Pointer(&decoder.data[0]))

	ret := C.opus_decoder_init(decoder.cDecoder, C.opus_int32(sampleRate), C.int(channels))
	if err := getErr(ret); err != nil {
		return nil, err
	}
	decoder.channels = channels

	return decoder, nil
}

func (d *Decoder) Decode(data []byte, frameSize int, fec bool) ([]int16, error) {
	var dataPtr *C.uchar
	if len(data) > 0 {
		dataPtr = (*C.uchar)(unsafe.Pointer(&data[0]))
	}
	dataLen := C.opus_int32(len(data))

	output := make([]int16, d.channels*frameSize)
	outputPtr := (*C.opus_int16)(unsafe.Pointer(&output[0]))

	var cFec C.int
	if fec {
		cFec = 1
	} else {
		cFec = 0
	}

	cRet := C.opus_decode(d.cDecoder, dataPtr, dataLen, outputPtr, C.int(frameSize), cFec)
	ret := int(cRet)

	if ret < 0 {
		return nil, getErr(cRet)
	}
	return output[:ret*d.channels], nil
}

func (d *Decoder) ResetState() {
	C.gopus_decoder_resetstate(d.cDecoder)
}

func CountFrames(data []byte) (int, error) {
	dataPtr := (*C.uchar)(unsafe.Pointer(&data[0]))
	cLen := C.opus_int32(len(data))

	cRet := C.opus_packet_get_nb_frames(dataPtr, cLen)
	if err := getErr(cRet); err != nil {
		return 0, err
	}
	return int(cRet), nil
}

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
