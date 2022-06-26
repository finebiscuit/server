package workspace

import "errors"

var (
	ErrNotFound         = errors.New("workspace not found")
	ErrPermissionDenied = errors.New("permission denied")
)
