package media

import (
	"errors"
	"strconv"

	fluentffmpeg "github.com/modfy/fluent-ffmpeg"
)

// Im using pointers for the fields so that I can see if a field is unset or not.

type ProbeResults struct {
	Video    Video
	Audio    []Audio
	Bitrate  *int64   // e.g. 5000000	// 5 Mbps	// Total bitrate of the media file
	Duration *float64 // e.g. 120.0 (seconds) // Total duration of the media file
	FileSize *int64   // e.g. 10485760 // 10 MB	// Total file size of the media file in bytes
}

type Video struct {
	Codec            *string  // e.g. "h264", "hevc"
	Bitrate          *int64   // e.g. 5000000	// 5 Mbps	// Total bitrate of the video stream
	Height           *int64   // e.g. 1080
	Width            *int64   // e.g. 1920
	Fps              *string  // e.g. 29.97
	Duration         *float64 // e.g. 120.0 (seconds)	// Total duration of the video stream
	ColorTransfer    *string  // e.g. "smpte2084", "arib-std-b67"	// Color transfer function used in the video
	BitsPerRawSample *int64   // e.g. 8, 10, 12	// Bits per raw sample (e.g. 8 for 8-bit video)
}

type Audio struct {
	Codec    *string // e.g. "aac", "mp3"
	Bitrate  *int64  // e.g. 128000	// 128 kb
	Channels *int64  // e.g. 2 (stereo)
}

func GetMediaValues(path string) (ProbeResults, error) {
	data, err := fluentffmpeg.Probe(path)
	if err != nil {
		return ProbeResults{}, err
	}

	return getMediaValues(data)
}

func getMediaValues(data map[string]interface{}) (ProbeResults, error) {
	formatRaw, ok := data["format"]
	if !ok {
		return ProbeResults{}, errors.New("format not found in probe data")
	}
	format, ok := formatRaw.(map[string]interface{})
	if !ok {
		return ProbeResults{}, errors.New("format is not a map")
	}

	streamsRaw, ok := data["streams"]
	if !ok {
		return ProbeResults{}, errors.New("streams not found in probe data")
	}

	streamsArray, ok := streamsRaw.([]interface{})
	if !ok {
		return ProbeResults{}, errors.New("streams is not an array")
	}
	if len(streamsArray) == 0 {
		return ProbeResults{}, errors.New("streams array is empty")
	}

	var video map[string]interface{}
	var audio []map[string]interface{}

	for _, stream := range streamsArray {
		streamMap, ok := stream.(map[string]interface{})
		if !ok {
			continue
		}

		codecType, ok := streamMap["codec_type"].(string)
		if !ok {
			continue
		}

		if codecType == "video" && video == nil {
			video = streamMap
		} else if codecType == "audio" {
			audio = append(audio, streamMap)
		}
	}

	var results ProbeResults
	getFormatValues(&results, format)
	getVideoValues(&results, video)
	getAudioValues(&results, audio)

	return results, nil
}

func getFormatValues(results *ProbeResults, format map[string]interface{}) {
	// duration
	if value, exists := format["duration"]; exists {
		if s, ok := value.(string); ok {
			if f, ok := strconv.ParseFloat(s, 64); ok != nil {
				results.Duration = &f
			}
		} else if f, ok := value.(float64); ok {
			results.Duration = &f
		}
	}

	// Bitrate
	if value, exists := format["bit_rate"]; exists {
		if s, ok := value.(string); ok {
			if f, ok := strconv.ParseInt(s, 10, 64); ok != nil {
				results.Bitrate = &f
			}
		} else if i, ok := value.(int64); ok {
			results.Bitrate = &i
		}
	}

	// File size
	if value, exists := format["size"]; exists {
		if s, ok := value.(string); ok {
			if f, ok := strconv.ParseInt(s, 10, 64); ok != nil {
				results.FileSize = &f
			}
		} else if i, ok := value.(int64); ok {
			results.FileSize = &i
		}
	}
}

func getVideoValues(results *ProbeResults, video map[string]interface{}) {
	// duration
	if value, exists := video["duration"]; exists {
		if s, ok := value.(string); ok {
			if f, ok := strconv.ParseFloat(s, 64); ok != nil {
				results.Video.Duration = &f
			}
		} else if f, ok := value.(float64); ok {
			results.Video.Duration = &f
		}
	}

	// Bitrate
	if value, exists := video["bit_rate"]; exists {
		if s, ok := value.(string); ok {
			if i, ok := strconv.ParseInt(s, 10, 64); ok != nil {
				results.Video.Bitrate = &i
			}
		} else if i, ok := value.(int64); ok {
			results.Video.Bitrate = &i
		}
	}

	// Codec
	if value, exists := video["codec_name"]; exists {
		if s, ok := value.(string); ok {
			results.Video.Codec = &s
		}
	}

	// Resolution
	if value, exists := video["height"]; exists {
		if s, ok := value.(string); ok {
			if i, ok := strconv.ParseInt(s, 10, 64); ok != nil {
				results.Video.Height = &i
			}
		} else if i, ok := value.(int64); ok {
			results.Video.Height = &i
		}
	}
	if value, exists := video["width"]; exists {
		if s, ok := value.(string); ok {
			if i, ok := strconv.ParseInt(s, 10, 64); ok != nil {
				results.Video.Width = &i
			}
		} else if i, ok := value.(int64); ok {
			results.Video.Width = &i
		}
	}

	// FPS
	if value, exists := video["r_frame_rate"]; exists {
		if s, ok := value.(string); ok {
			results.Video.Fps = &s
		}
	}

	// Color Transfer
	if value, exists := video["color_transfer"]; exists {
		if s, ok := value.(string); ok {
			results.Video.ColorTransfer = &s
		}
	}

	// Bits per raw sample
	if value, exists := video["bits_per_raw_sample"]; exists {
		if s, ok := value.(string); ok {
			if i, ok := strconv.ParseInt(s, 10, 64); ok != nil {
				results.Video.BitsPerRawSample = &i
			}
		} else if i, ok := value.(int64); ok {
			results.Video.BitsPerRawSample = &i
		}
	}
}

func getAudioValues(results *ProbeResults, audio []map[string]interface{}) {

	for _, stream := range audio {
		var streamResult Audio

		// Codec
		if value, exists := stream["codec_name"]; exists {
			if s, ok := value.(string); ok {
				streamResult.Codec = &s
			}
		}

		// Bitrate
		if value, exists := stream["bit_rate"]; exists {
			if s, ok := value.(string); ok {
				if i, ok := strconv.ParseInt(s, 10, 64); ok != nil {
					streamResult.Bitrate = &i
				}
			} else if i, ok := value.(int64); ok {
				streamResult.Bitrate = &i
			}
		}

		// Channels
		if value, exists := stream["channels"]; exists {
			if s, ok := value.(string); ok {
				if i, ok := strconv.ParseInt(s, 10, 64); ok != nil {
					streamResult.Channels = &i
				}
			} else if i, ok := value.(int64); ok {
				streamResult.Channels = &i
			}
		}

		results.Audio = append(results.Audio, streamResult)
	}
}
