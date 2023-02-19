package downloader

const (
	FormatAudio = "audio"
	FormatVideo = "video"
)

type Format string

func (f Format) Valid() bool {
	return f == FormatAudio || f == FormatVideo
}
