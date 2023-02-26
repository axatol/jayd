package downloader

import "fmt"

func itemFilename(item QueueItem) string {
	var selectedFormat Format
	for _, format := range item.Data.Formats {
		if item.FormatID == format.FormatID {
			selectedFormat = format
			break
		}
	}

	ext := selectedFormat.VideoExt
	if selectedFormat.VideoExt == "none" {
		ext = selectedFormat.AudioExt
	}

	return fmt.Sprintf("%s_%s.%s", item.Data.ID, item.FormatID, ext)
}
