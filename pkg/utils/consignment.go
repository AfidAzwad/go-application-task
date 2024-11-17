package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// GenerateConsignmentID generates a unique consignment ID based on the provided format
func GenerateConsignmentID(cityCode string) (string, error) {
	currentDate := time.Now().Format("060102") // "YYMMDD" format

	// Generate a random alphanumeric identifier of 4 characters
	identifier := generateRandomString(4)

	consignmentID := fmt.Sprintf("CID%s%s%s", currentDate, cityCode, identifier)
	return consignmentID, nil
}

// generateRandomString generates a random alphanumeric string of the given length
func generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := 0; i < length; i++ {
		sb.WriteByte(charset[r.Intn(len(charset))])
	}

	return sb.String()
}
