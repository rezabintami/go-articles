package helpers

import (
	"strconv"
	"time"
)

func GenerateReceipt() string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	receipt := "INV" + timestamp
	return receipt
}

func GenerateManifest() string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	receipt := "MNF" + timestamp
	return receipt
}