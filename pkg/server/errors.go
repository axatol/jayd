package server

import "fmt"

var (
	err_MissingFormat = fmt.Errorf("missing format")
	err_MissingTarget = fmt.Errorf("missing target")
	err_MissingKey    = fmt.Errorf("missing key")
	err_InvalidFormat = fmt.Errorf("invalid format")
	err_InvalidTarget = fmt.Errorf("invalid target")
	err_FailedRequest = fmt.Errorf("failed request")
	err_Unauthorised  = fmt.Errorf("unauthorised")
	err_FetchMetadata = fmt.Errorf("failed to fetch metadata")
	err_NotFound      = fmt.Errorf("item was not found")
	err_GenericError  = fmt.Errorf("an error occurred")
)
