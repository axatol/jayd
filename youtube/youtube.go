package youtube

type ListResponse[T any] struct {
	Kind  string `json:"kind"`
	ETag  string `json:"etag"`
	Items []T    `json:"items"`
}

type Client struct{}
