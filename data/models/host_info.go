package models

type HostInfo struct {
	OperatingSystem string `json:"operating_system"`
	Architecture    string `json:"architecture"`
	CPUnNumber      int    `json:"cpu_number"`
	CurrentUnixNano int64  `json:"current_unix_nano"`
}
