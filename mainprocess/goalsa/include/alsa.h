/*

    this contain the api headers.

*/
#define ALSA_PCM_NEW_HW_PARAMS_API
#include <alsa/asoundlib.h>
#include <stdlib.h>

struct goalsa_error {
	int error_code;
    char error_string[1024];
};

struct goalsa_device {
    unsigned int wpm;
    unsigned int frequency;
    char name[255];
    snd_pcm_t *handle;
    snd_pcm_uframes_t frames; // 32
    unsigned int sampling_rate; // 44100
    int number_channels; // 2
    int periods;
    int dir; // exact_rate - rate (-1, 0, 1).
    snd_pcm_format_t format; // SND_PCM_FORMAT_S16_LE | SND_PCM_FORMAT_FLOAT
    snd_pcm_access_t access; // SND_PCM_ACCESS_RW_INTERLEAVED
    snd_pcm_uframes_t frames_per_period;
    snd_pcm_uframes_t buffer_size;
    snd_pcm_hw_params_t *hardware_params;
    // 1 frame = goalsa_device.number_channels * goalsa_nbytes_in_format((*goalsa_device).format)
};

struct goalsa_device_error {
	struct goalsa_device *device;
	struct goalsa_error *error;
};

struct goalsa_info {
    char *version;
    unsigned int stream_type_count;
    unsigned int access_type_count;
    unsigned int format_count;
    unsigned int subformat_count;
    unsigned int state_count;
};

struct goalsa_value_name {
    int value;
    const char *name;
};

struct goalsa_value_name_description {
    int value;
    const char *name;
    const char *description;
};
