package goalsa

import (
	"math"
	"time"
	"unsafe"

	"github.com/pkg/errors"
)

/*

// Use the newer ALSA API
#define ALSA_PCM_NEW_HW_PARAMS_API
#cgo LDFLAGS: -lasound
//cgo CFLAGS: -lpcm -Iinclude
#cgo CFLAGS: -Iinclude
#include <alsa/asoundlib.h>
#include <stdlib.h>
#include "alsa.h"
#include "math.h"

struct goalsa_device *goalsa_no_device = (struct goalsa_device *)0;
snd_pcm_hw_params_t *goalsa_no_hw_params = (snd_pcm_hw_params_t *)0;

int goalsa_nbytes_in_format(snd_pcm_format_t format) {
	if (format == SND_PCM_FORMAT_FLOAT) {
		return 4;
	}
    if (format == SND_PCM_FORMAT_U32_LE) {
        return 4;
    }
    if (format == SND_PCM_FORMAT_U32_BE) {
        return 4;
    }
    if (format == SND_PCM_FORMAT_FLOAT64_LE) {
        return 8;
    }
    if (format == SND_PCM_FORMAT_S24_3BE) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_S20_3BE) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_U8) {
        return 1;
    }
    if (format == SND_PCM_FORMAT_S16_LE) {
        return 2;
    }
    if (format == SND_PCM_FORMAT_U24_3BE) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_G723_40) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_DSD_U8) {
        return 1;
    }
    if (format == SND_PCM_FORMAT_U16_LE) {
        return 2;
    }
    if (format == SND_PCM_FORMAT_U24_3LE) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_S24_3LE) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_FLOAT_LE) {
        return 4;
    }
    if (format == SND_PCM_FORMAT_S32_LE) {
        return 4;
    }
    if (format == SND_PCM_FORMAT_FLOAT_BE) {
        return 4;
    }
    if (format == SND_PCM_FORMAT_S18_3BE) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_S16_BE) {
        return 2;
    }
    if (format == SND_PCM_FORMAT_S24_BE) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_G723_24) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_DSD_U32_LE) {
        return 4;
    }
    if (format == SND_PCM_FORMAT_S20_3LE) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_U20_3LE) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_DSD_U16_BE) {
        return 2;
    }
    if (format == SND_PCM_FORMAT_S24_LE) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_U18_3BE) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_G723_40_1B) {
        return 1;
    }
    if (format == SND_PCM_FORMAT_U16_BE) {
        return 2;
    }
    if (format == SND_PCM_FORMAT_FLOAT64_BE) {
        return 8;
    }
    if (format == SND_PCM_FORMAT_S18_3LE) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_U18_3LE) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_DSD_U32_BE) {
        return 4;
    }
    if (format == SND_PCM_FORMAT_U24_LE) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_S32_BE) {
        return 4;
    }
    if (format == SND_PCM_FORMAT_U20_3BE) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_G723_24_1B) {
        return 3;
    }
    if (format == SND_PCM_FORMAT_DSD_U16_LE) {
        return 2;
    }
    if (format == SND_PCM_FORMAT_U24_BE) {
        return 3;
    }
    // SND_PCM_FORMAT_GSM, GSM
    // SND_PCM_FORMAT_IMA_ADPCM, Ima-ADPCM
    // SND_PCM_FORMAT_SPECIAL, Special
    // SND_PCM_FORMAT_A_LAW, A-Law
    // SND_PCM_FORMAT_IEC958_SUBFRAME_BE, IEC-958 Big Endian
    // SND_PCM_FORMAT_MU_LAW, Mu-Law
    // SND_PCM_FORMAT_MPEG, MPEG
    // SND_PCM_FORMAT_IEC958_SUBFRAME_LE, IEC-958 Little Endian
    return 1;
}

struct goalsa_device *GoAlsaNewPlayer(unsigned int frequency, unsigned int wpm) {
	struct goalsa_device *v;
	v = (struct goalsa_device *) malloc(sizeof(struct goalsa_device));
	strcpy((char *)(*v).name, "default");
	(*v).frequency = frequency;
	(*v).wpm =  wpm;
    (*v).frames = 32; // 1 frame = .number_channels * goalsa_nbytes_in_format((*v).format)
    (*v).sampling_rate = 44100;
    (*v).number_channels = 2;
	(*v).frames_per_period = 4410;
	(*v).format = SND_PCM_FORMAT_S16_LE;
    (*v).access = SND_PCM_ACCESS_RW_INTERLEAVED;
	(*v).hardware_params = goalsa_no_hw_params;
    (*v).buffer_size = 0;
    return v;
}

struct goalsa_device_error *newDeviceError() {
	struct goalsa_device_error *v;
	v = (struct goalsa_device_error *) malloc(sizeof(struct goalsa_device_error));
    (*v).device = goalsa_no_device;
    //(*v).error = newError();
	(*v).error = malloc(sizeof(struct goalsa_error));
    return v;
}

const char *GoAlsaDeviceErrorGetErrorString(struct goalsa_device_error *de) {
    return (*(*de).error).error_string;
    //return (*de).error->error_string;
}

struct goalsa_device_error *GoAlsaOpenPlayer(struct goalsa_device *device) {
	int rc;
	int dir;
	unsigned int val;
	struct goalsa_device_error *return_value = newDeviceError();

	//.open PCM device for playback.
	rc = snd_pcm_open(&(*device).handle, (*device).name, SND_PCM_STREAM_PLAYBACK, 0);
	if (rc < 0) {
		(*(*return_value).error).error_code = rc;
		sprintf((*(*return_value).error).error_string, "unable to open pcm device: %s", snd_strerror(rc));
		return return_value;
	}

	// Allocate a hardware parameters object.
	snd_pcm_hw_params_alloca(&(*device).hardware_params);

	// Fill it in with default values.
	snd_pcm_hw_params_any((*device).handle, (*device).hardware_params);

	// Set the desired hardware parameters.

	// mode
	snd_pcm_hw_params_set_access((*device).handle, (*device).hardware_params, (*device).access);

	// format
	snd_pcm_hw_params_set_format((*device).handle, (*device).hardware_params, (*device).format);

	// sampling rate
	snd_pcm_hw_params_set_rate_near((*device).handle, (*device).hardware_params, &(*device).sampling_rate, &dir);

	// # of channels
	snd_pcm_hw_params_set_channels((*device).handle, (*device).hardware_params, (*device).number_channels);

	unsigned int samplesPerSecond = (*device).sampling_rate / (*device).frequency;
	unsigned int nElementsPerMinute = 50 * (*device).wpm; // 50 == # elements in "paris"
	float nElementsPerSecond = floor((float)nElementsPerMinute / (float)60);
	unsigned int samplesPerElement = (*device).sampling_rate / nElementsPerSecond;
	(*device).frames_per_period = samplesPerElement;

	// Set period size in frames.
	snd_pcm_hw_params_set_period_size_near((*device).handle, (*device).hardware_params, &(*device).frames_per_period, &dir);

	// Set buffer size (in frames). The resulting latency is given by
	// latency = periodsize * periods / (rate * bytes_per_frame)
	int frame_size = (*device).number_channels * goalsa_nbytes_in_format((*device).format);
	(*device).buffer_size = (*device).frames_per_period * frame_size;
    snd_pcm_hw_params_set_buffer_size_near((*device).handle, (*device).hardware_params, (snd_pcm_uframes_t *)&(*device).buffer_size);

	// Write the parameters to the driver
	rc = snd_pcm_hw_params((*device).handle, (*device).hardware_params);
	if (rc < 0) {
		(*(*return_value).error).error_code = rc;
		sprintf((*(*return_value).error).error_string, "unable to set hw parameters: %s", snd_strerror(rc));
		return return_value;
	}

	(*return_value).device = device;
	return return_value;
}

void GoAlsaCloseDevice(struct goalsa_device *device) {
  snd_pcm_drain((*device).handle);
  snd_pcm_close((*device).handle);
}

void GoAlsaFreeHWParams(struct goalsa_device *device) {
	// this always pukes so dont use it
	snd_pcm_hw_params_free((*device).hardware_params);
}

void GoAlsaDeconstructDevice(struct goalsa_device *device) {
    if (device != goalsa_no_device) {
        free(device);
        device = goalsa_no_device;
    }
}

void free_goalsa_info(struct goalsa_info *v) {
	free((*v).version);
	free(v);
}

void free_goalsa_value_name(struct goalsa_value_name *v) {
	free(v);
}

void free_goalsa_value_name_description(struct goalsa_value_name_description *v) {
	free(v);
}

struct goalsa_info *goalsa_get_info() {
	struct goalsa_info *v;
	v = (struct goalsa_info *) malloc(sizeof(struct goalsa_info));
	(*v).version = malloc(sizeof(SND_LIB_VERSION_STR));
	strcpy((*v).version, SND_LIB_VERSION_STR);
    (*v).stream_type_count = SND_PCM_STREAM_LAST;
    (*v).access_type_count = SND_PCM_ACCESS_LAST;
    (*v).format_count = SND_PCM_FORMAT_LAST;
    (*v).subformat_count = SND_PCM_SUBFORMAT_LAST;
    (*v).state_count = SND_PCM_STATE_LAST;
	return v;
}

struct goalsa_value_name *goalsa_get_stream_type_name(unsigned int index) {
	struct goalsa_value_name *v;
	v = (struct goalsa_value_name *) malloc(sizeof(struct goalsa_value_name));
	if (index <= SND_PCM_STREAM_LAST) {
		(*v).value = (int)index;
		(*v).name = snd_pcm_stream_name((snd_pcm_stream_t)index);
	} else {
		(*v).value = -1;
	}
	return v;
}

struct goalsa_value_name *goalsa_get_access_type_name(unsigned int index) {
	struct goalsa_value_name *v;
	v = (struct goalsa_value_name *) malloc(sizeof(struct goalsa_value_name));
	if (index <= SND_PCM_ACCESS_LAST) {
		(*v).value = (int)index;
		(*v).name = snd_pcm_access_name((snd_pcm_access_t)index);
	} else {
		(*v).value = -1;
	}
	return v;
}

struct goalsa_value_name_description *goalsa_get_format_name_description(unsigned int index) {
	struct goalsa_value_name_description *v;
	v = (struct goalsa_value_name_description *) malloc(sizeof(struct goalsa_value_name_description));
	if (index <= SND_PCM_FORMAT_LAST) {
		(*v).value = (int)index;
		(*v).name = snd_pcm_format_name((snd_pcm_format_t)index);
		(*v).description = snd_pcm_format_description((snd_pcm_format_t)index);
	} else {
		(*v).value = -1;
	}
	return v;
}

struct goalsa_value_name_description *goalsa_get_subformat_name_description(unsigned int index) {
	struct goalsa_value_name_description *v;
	v = (struct goalsa_value_name_description *) malloc(sizeof(struct goalsa_value_name_description));
	if (index <= SND_PCM_SUBFORMAT_LAST) {
		(*v).value = (int)index;
		(*v).name, snd_pcm_subformat_name((snd_pcm_subformat_t)index);
		(*v).description, snd_pcm_subformat_description((snd_pcm_subformat_t)index);
	} else {
		(*v).value = -1;
	}
	return v;
}

struct goalsa_value_name *goalsa_get_state_name(unsigned int index) {
	struct goalsa_value_name *v;
	v = (struct goalsa_value_name *) malloc(sizeof(struct goalsa_value_name));
	if (index <= SND_PCM_STATE_LAST) {
		(*v).value = (int)index;
		(*v).name, snd_pcm_state_name((snd_pcm_state_t)index);
	} else {
		(*v).value = -1;
	}
	return v;
}

void GoAlsaMakeNote(struct goalsa_device *device, unsigned int n_elements) {
	unsigned int n_samples = n_elements * (*device).frames_per_period;
	unsigned int frame_size = (*device).number_channels * goalsa_nbytes_in_format((*device).format);
	unsigned int mallocSize  = n_samples * frame_size;

	unsigned char *data = (unsigned char *)malloc(mallocSize);
	// each sample is a frame.
	for(int i = 0; i < n_samples; i++) {
        short s1 = (i % 128) * 100 - 10000;
        short s2 = (i % 256) * 100 - 10000;
        data[4*i] = (unsigned char)s1;
        data[4*i+1] = s1 >> 8;
        data[4*i+2] = (unsigned char)s2;
        data[4*i+3] = s2 >> 8;
	}

	int rc = snd_pcm_writei((*device).handle, data, n_samples);
	while (rc < 0) {
		// if (rc == -EPIPE) then underrun
		// if (rc < 0) some other error.
		snd_pcm_prepare((*device).handle);
		rc = snd_pcm_writei((*device).handle, data, n_samples);
	}
}

*/
import "C"

func getAlsaInfo() *AlsaInfo {
	var v AlsaInfo
	info := C.goalsa_get_info()
	defer C.free_goalsa_info(info)
	v.Version = C.GoString(info.version)
	v.StreamTypes = make(map[int]string)
	for i := 0; i <= int(info.stream_type_count); i++ {
		r := C.goalsa_get_stream_type_name(C.uint(i))
		if r.value >= 0 {
			v.StreamTypes[i] = C.GoString(r.name)
		}
		C.free_goalsa_value_name(r)
	}
	v.AccessTypes = make(map[int]string)
	for i := 0; i <= int(info.access_type_count); i++ {
		r := C.goalsa_get_access_type_name(C.uint(i))
		if r.value > 0 {
			v.AccessTypes[i] = C.GoString(r.name)
		}
		C.free_goalsa_value_name(r)
	}
	v.Formats = make(map[int]struct{ Name, Description string })
	for i := 0; i <= int(info.format_count); i++ {
		r := C.goalsa_get_format_name_description(C.uint(i))
		if r.value > 0 {
			v.Formats[i] = struct{ Name, Description string }{Name: C.GoString(r.name), Description: C.GoString(r.description)}
		}
		C.free_goalsa_value_name_description(r)
	}
	v.SubFormats = make(map[int]struct{ Name, Description string })
	for i := 0; i <= int(info.subformat_count); i++ {
		r := C.goalsa_get_subformat_name_description(C.uint(i))
		if r.value > 0 {
			v.SubFormats[i] = struct{ Name, Description string }{Name: C.GoString(r.name), Description: C.GoString(r.description)}
		}
		C.free_goalsa_value_name_description(r)
	}
	v.States = make(map[int]string)
	for i := 0; i <= int(info.state_count); i++ {
		r := C.goalsa_get_state_name(C.uint(i))
		if r.value >= 0 {
			v.States[i] = C.GoString(r.name)
		}
		C.free_goalsa_value_name(r)
	}
	return &v
}

// Player is an audio device.
// It contains a pointer to a C struct goalsa_device.
// The C struct goalsa_device is the output device information.
type Player struct {
	device *C.struct_goalsa_device
}

// newPlayer constructs a new audio device with a byte buffer.
func newPlayer(frequency, wpm uint64) *Player {
	var v Player
	v.device = C.GoAlsaNewPlayer(C.uint(frequency), C.uint(wpm))
	return &v
}

//open opens the device.
// This sets the device to a valid device.
// Returns the error or nil.
func (player *Player) open() (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(player *Player)open()")
		}
	}()
	returnStruct := C.GoAlsaOpenPlayer(player.device)
	defer freeDeviceError(returnStruct)
	if (*returnStruct.error).error_code < 0 {
		cmessage := C.GoAlsaDeviceErrorGetErrorString(returnStruct)
		defer C.free(unsafe.Pointer(cmessage))
		message := C.GoString(cmessage)
		err = errors.New(message)
		return
	}
	player.device = returnStruct.device
	return
}

// close flushes the device.
// Frees up allocated memory used by device.
// open must be called before using Play again.
func (player *Player) close() {
	C.GoAlsaCloseDevice(player.device)
	C.GoAlsaDeconstructDevice(player.device)
	player.device = nil
	// the following line pukes.
	///player.device = C.goalsa_no_device
}

/*

If you really want to do it with that function, generate a waveform in a buffer.
A triangle-shaped wave may not sound too awful and should be simple enough to generate.

The base "la" (A) is 440Hz, that is, 440 cycles of the waveform of your choice per second.
The other notes can be obtained by multiplying/dividing by 2^(1/12) ( 1.05946309 )
 for each half tone above/below this base frequency.

If the device frequency is, say, 44100 Hz, and you want to play the base "la",
 each period of your waveform should occupy 44100 / 440 or about 100 samples.

Pay attention to the sample width and the number of channels the device is configured for, too.

Explanation:
 there are 12 half tones in an octave,
 and an octave is exactly half (lower pitched) or double (higher pitched) the frequency.
 Once you have multiplied 12 times by 2^(1/12), you have multiplied by 2,
  so each half-tone is at a factor of 2^(1/12) above the previous one.

*/

// ditDahToFF converts ".- -."  to bytes
func (player *Player) ditDahToFF(ditdah string) (err error) {
	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(player *Player) ditDahToF(ditdah string, wpm, frequency uint64)")
		}
	}()
	// word == "paris" == 50 elements.
	//samplesPerPeriod := float64(player.device.sampling_rate) / float64(frequency)
	//samplesPerElement := float64(player.device.sampling_rate) / nElementsPerSecond
	nElementsPerMinute := float64(50) * float64((*player.device).wpm)
	nElementsPerSecond := math.Floor(nElementsPerMinute / float64(60))
	secondsPerElement := float64(1) / nElementsPerSecond

	var elementCount uint64
	var pauseCount float64
	for _, r := range ditdah {
		switch r {
		case ' ':
			// space is used as a character separator.
			// char separator : 3 elements of silence.
			// 2 pauses added to the pause following the last dit or dah.
			elementCount = 0
			pauseCount = 2
		case '\t':
			// tab is used as a word separator.
			// word separator : 7 elements of silence.
			// 6 pauses added to the pause following the last dit or dah.
			elementCount = 0
			pauseCount = 6
		case '.':
			// period is used as a dit.
			// dit : 1 element of sound followed by 1 element of silence
			elementCount = 1
			pauseCount = 1
		case '-':
			// dash is used as a dah.
			// dah : 3 elements of sound followed by 1 element of silence.
			elementCount = 3
			pauseCount = 1
		case 's':
			// 1 second dah
			elementCountf := math.Floor(nElementsPerSecond)
			elementCount = uint64(elementCountf)
			pauseCount = 0
		case 'm':
			// 30 second dah
			elementCountf := math.Floor(nElementsPerMinute / float64(2.0))
			elementCount = uint64(elementCountf)
			pauseCount = 0
		}
		// buffer -> period -> frames:32 -> samples:2 -> float64
		if elementCount > 0 {
			// sound
			elementCountF := float64(elementCount)
			// timeout after the sound and the pause.
			pauseFor := (secondsPerElement * elementCountF)
			dur := uint64(float64(time.Second) * pauseFor)
			timeout := time.After(time.Duration(dur))

			// play the sound
			C.GoAlsaMakeNote(player.device, C.uint(elementCount))
			// wait
			select {
			case <-timeout:
			}

		}
		if pauseCount > 0 {
			pauseFor := secondsPerElement * pauseCount
			//fmt.Printf("pausing for %f seconds\n", pauseFor)
			dur := uint64(float64(time.Second) * pauseFor)
			timeout := time.After(time.Duration(dur))
			select {
			case <-timeout:
			}
		}
	}
	return
}

func freeDeviceError(returnstruct *C.struct_goalsa_device_error) {
	C.free(unsafe.Pointer(returnstruct.error))
	C.free(unsafe.Pointer(returnstruct))
}
