package domain

import (
	"github.com/pkg/errors"
)

var (
	ErrScriptNotFound    = errors.New("script not found")
	ErrDuplicateScriptId = errors.New("duplicate script id")
)
