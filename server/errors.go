package server

import "fmt"

var (
	err_MissingTarget = fmt.Errorf("missing target")
	err_InvalidFormat = fmt.Errorf("invalid format")
	err_InvalidTarget = fmt.Errorf("invalid target")
	err_FailedRequest = fmt.Errorf("failed request")
	err_Unauthorised  = fmt.Errorf("unauthorised")
	err_FetchMetadata = fmt.Errorf("failed to fetch metadata")
)
