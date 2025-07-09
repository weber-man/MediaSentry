package media

import (
	"log"
)

func Checks(path string) {

	log.Println("Running checks for", path)
	probe, err := GetMediaValues(path)
	if err != nil {
		log.Println("Error getting media values:", err)
		return
	}

	isHdr := isHdr(probe)
	log.Println("is HDR:", isHdr)
}

func isHdr(probe ProbeResults) bool {

	// For HDR detection ColorTransfer is checked.
	// BitsPerRawSample is only used if ColorTransfer is not available.
	if probe.Video.ColorTransfer != nil {
		switch *probe.Video.ColorTransfer {
		case "smpte2084", "arib-std-b67":
			return true
		default:
			return false
		}
	}

	if probe.Video.BitsPerRawSample != nil {
		return *probe.Video.BitsPerRawSample >= 10
	}
	return false
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
