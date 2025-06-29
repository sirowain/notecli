package localdb

import (
	"errors"
)

// ErrBucketNotFound is returned when the specified bucket does not exist
var ErrBucketNotFound = errors.New("bucket not found")
