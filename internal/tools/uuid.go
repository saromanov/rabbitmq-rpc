package tools

import "github.com/google/uuid"

// GenerateUUID provides generating of uuid
func GenerateUUID() string {
	return uuid.NewString()
}