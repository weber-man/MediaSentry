package config

import (
	"os"
	"strconv"
	"strings"
)

func MediaFolder() string {
	var mediaFolder string = "/media"
	if os.Getenv("MEDIA_FOLDER") != "" {
		mediaFolder = os.Getenv("MEDIA_FOLDER")
	}
	return mediaFolder
}

func FileExtensions() []string {
	var extensions []string = []string{".mp4", ".mkv", ".avi", ".mov", ".webm", ".flv"}
	if os.Getenv("FILE_EXTENSIONS") != "" {
		extensions = strings.Split(os.Getenv("FILE_EXTENSIONS"), ",")
	}
	return extensions
}

func CheckAllowHDR() bool {
	var allowHDR bool = true
	if os.Getenv("CHECK_ALLOW_HDR") != "" {
		if os.Getenv("CHECK_ALLOW_HDR") == "false" {
			allowHDR = false
		}
	}
	return allowHDR
}

func CheckMaxVideoFps() int {
	if os.Getenv("CHECK_MAX_VIDEO_FPS") != "" {
		if fps, err := strconv.Atoi(os.Getenv("CHECK_MAX_VIDEO_FPS")); err == nil {
			return fps
		}
	}
	return 60
}

func CheckMaxVideoWidth() int {
	if os.Getenv("CHECK_MAX_VIDEO_WIDTH") != "" {
		if width, err := strconv.Atoi(os.Getenv("CHECK_MAX_VIDEO_WIDTH")); err == nil {
			return width
		}
	}
	return 3840
}

func CheckMaxVideoHeight() int {
	if os.Getenv("CHECK_MAX_VIDEO_HEIGHT") != "" {
		if height, err := strconv.Atoi(os.Getenv("CHECK_MAX_VIDEO_HEIGHT")); err == nil {
			return height
		}
	}
	return 2160
}

func CheckMaxVideoBitrate() int64 {
	if os.Getenv("CHECK_MAX_VIDEO_BITRATE") != "" {
		if bitrate, err := strconv.ParseInt(os.Getenv("CHECK_MAX_VIDEO_BITRATE"), 10, 64); err == nil {
			return bitrate
		}
	}
	return 5000000 // Default to 5 Mbps
}

func CheckAllowedCodecs() []string {
	var allowedCodecs []string = []string{"h264", "hevc", "av1"}
	if os.Getenv("CHECK_ALLOWED_CODECS") != "" {
		allowedCodecs = strings.Split(os.Getenv("CHECK_ALLOWED_CODECS"), ",")
	}
	return allowedCodecs
}

func CheckMaxFileSize() int64 {
	if os.Getenv("CHECK_FILE_SIZE_LIMIT") != "" {
		if size, err := strconv.ParseInt(os.Getenv("CHECK_FILE_SIZE_LIMIT"), 10, 64); err == nil {
			return size
		}
	}
	return 21474836480 // Default to 20 GB
}
