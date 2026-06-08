// Package uuid provides UUID generation.
package uuid

import "github.com/google/uuid"

func Generate() string {
	return uuid.New().String()
}
