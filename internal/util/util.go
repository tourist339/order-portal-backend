package util

import (
	"github.com/google/uuid"
	"strconv"
	"time"
)

func GenerateUniqueID(prefix string) string {
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)

	// Obtain a unique code using UUID
	uniqueID := uuid.New().ID()

	// Combine the time-based component and the UUID to create a globally unique ID
	ID := currentTimestamp + int64(uniqueID)
	return prefix + "-" + strconv.FormatInt(ID, 10)
}
