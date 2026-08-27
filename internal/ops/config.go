package ops

import "os"

func DataDir() string {
	if v := os.Getenv("RAILGUARD_DATA"); v != "" {
		return v
	}
	return "./railguard-data"
}
