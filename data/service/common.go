package service

import (
	"filecoin-data-provider/models"
	"runtime"
	"time"
)

func GetHostInfo() *models.HostInfo {
	info := models.HostInfo{
		OperatingSystem: runtime.GOOS,
		Architecture:    runtime.GOARCH,
		CPUnNumber:      runtime.NumCPU(),
		CurrentUnixNano: time.Now().UnixNano(),
	}

	return &info
}
