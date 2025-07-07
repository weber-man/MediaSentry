package main

import (
	"log"
	fluentffmpeg "github.com/modfy/fluent-ffmpeg"
)

func checks(path string) {

	log.Println("Running checks for", path)
	isHdr, err := isHdr(path)
	if err != nil {
		return
	}
	log.Println("is HDR:", isHdr)
}

func isHdr(path string) (bool, error) {
	data, err := fluentffmpeg.Probe(path)
	if err != nil {
		return false, err
	}
	streams := data["streams"].([]interface{})
	var color interface{}
	var bps interface{}
	
	for _, stream := range streams {
		streamMap := stream.(map[string]interface{})
		if codecType, ok := streamMap["codec_type"]; ok && codecType == "video" {
			color = streamMap["color_transfer"]
			bps = streamMap["bits_per_raw_sample"]
			break
		}
	}

	if (color != nil && (color == "smpte2084" || color == "arib-std-b67") ) {
		return true, nil
	}
	if bps != nil {
		if bpsInt, ok := bps.(int); ok && bpsInt >= 10 {
			return true, nil
		}
		if bpsFloat, ok := bps.(float64); ok && bpsFloat >= 10 {
			return true, nil
		}
	}
	return false, nil
}

func videoCodecCheck(path string) bool {
	return true
}

func fileSizeCheck(path string) bool {
	return true
}

func videoResolutionCheck(path string) bool {
	return true
}

func videoBitrateCheck(path string) bool {
	return true
}

func videoFpsCheck(path string) bool {
	return true
}

func audioChannelsCheck(path string) bool {
	return true
}

func audioSampleRateCheck(path string) bool {
	return true
}

func audioBitrateCheck(path string) bool {
	return true
}

func audioCodecCheck(path string) bool {
	return true
}
