package filemanager

import (
	"github.com/google/uuid"
)

func NewFileName(filetype string) string {
	return uuid.NewString() + filetype
}