package datefile

import (
	"fmt"
	"os"
	"time"
)

func OutsideInterval(filePath string, interval time.Duration) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Unable to read previous date in file :", err)
		return true
	}
	previousCheck, err := time.Parse(time.DateTime, string(data))
	if err != nil {
		fmt.Println("Unable to parse previous date in file :", err)
		return true
	}
	return time.Now().Sub(previousCheck) > interval
}

func Write(filePath string) {
	if err := os.WriteFile(filePath, time.Now().AppendFormat(nil, time.DateTime), 0644); err != nil {
		fmt.Println("Unable to write date in file :", err)
	}
}
