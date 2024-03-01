// Package sudo provides error messages for the sudo package.
package sudo

// messages struct holds error messages for various scenarios in the sudo package.
type messages struct {
	FailedRedisFetch string
	FailedRedisSet   string
	FailedMarshal    string
	NotFound         string
	InvalidToken     string
	InvalidCode      string
	ExceedTryCount   string
	ExpiredCode      string

	VerifyStarted string
	Unknown       string
}

// Messages is an instance of the messages struct with predefined error messages.
var Messages = messages{
	FailedRedisFetch: "sudo_redis_fetch_failed",
	FailedRedisSet:   "sudo_redis_set_failed",
	FailedMarshal:    "sudo_json_marshal_failed",
	NotFound:         "sudo_not_found",
	InvalidToken:     "sudo_invalid_token",
	InvalidCode:      "sudo_invalid_code",
	ExpiredCode:      "sudo_expired_code",
	ExceedTryCount:   "sudo_exceed_try_count",
	Unknown:          "sudo_unknown",
	VerifyStarted:    "sudo_verify_started",
}
