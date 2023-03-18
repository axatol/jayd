package downloader

import (
	"fmt"
)

func itemFilename(info InfoJSON) string {
	var selectedFormat Format
	for _, format := range info.Formats {
		if info.FormatID == format.FormatID {
			selectedFormat = format
			break
		}
	}

	ext := selectedFormat.VideoExt
	if selectedFormat.VideoExt == "none" {
		ext = selectedFormat.AudioExt
	}

	return fmt.Sprintf("%s_%s.%s", info.VideoID, info.FormatID, ext)
}
