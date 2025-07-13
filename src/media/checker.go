package media

import (
	"log"
	"slices"
	"strconv"
	"strings"
)

func Checks(path string) {

	log.Println("Running checks for", path)
	probe, err := GetMediaValues(path)
	if err != nil {
		log.Println("Error getting media values:", err)
		return
	}

	if isHdr, ok := isHdr(probe); ok {
		log.Println("HDR check result:", isHdr)

	} else {
		log.Println("HDR check not applicable")
	}
}

func isHdr(probe ProbeResults) (result bool, ok bool) {

	// For HDR detection ColorTransfer is checked.
	// BitsPerRawSample is only used if ColorTransfer is not available.
	if probe.Video.ColorTransfer != nil {
		switch *probe.Video.ColorTransfer {
		case "smpte2084", "arib-std-b67":
			return true, true
		default:
			return false, true
		}
	}

	if probe.Video.BitsPerRawSample != nil {
		return *probe.Video.BitsPerRawSample >= 10, true
	}
	return false, false
}

func isCodecAllowed(probe ProbeResults, allowdCodecs []string) (result bool, ok bool) {
	if probe.Video.Codec != nil {
		return slices.Contains(allowdCodecs, *probe.Video.Codec), true
	}
	return false, false
}

func isFileSizeWithinLimit(probe ProbeResults, maxFileSizeInBytes int64) (result bool, ok bool) {
	if probe.FileSize != nil {
		return *probe.FileSize <= maxFileSizeInBytes, true
	}
	return false, false
}

func isVideoResolutionExceeded(probe ProbeResults, width int64, height int64) (result bool, ok bool) {
	if probe.Video.Width != nil && probe.Video.Height != nil {
		return *probe.Video.Width > width || *probe.Video.Height >
			height, true
	}
	return false, false
}

func isBitrateOverThreshold(probe ProbeResults, bitrate int64) (result bool, ok bool) {
	if probe.Video.Bitrate != nil {
		return *probe.Video.Bitrate > bitrate, true
	}
	return false, false
}

func isFpsWithinLimit(probe ProbeResults, maxFps float64) (result bool, ok bool) {
	if probe.Video.Fps != nil {
		strings := strings.Split(*probe.Video.Fps, "/")
		if len(strings) != 2 {
			return false, false
		}
		numerator, err1 := strconv.Atoi(strings[0])
		denominator, err2 := strconv.Atoi(strings[1])
		if err1 != nil || err2 != nil {
			return false, false
		}
		if denominator == 0 {
			return false, false
		}
		fps := float64(numerator) / float64(denominator)
		return fps > maxFps, true
	}
	return false, false
}
