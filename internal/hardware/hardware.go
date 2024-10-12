package hardware

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

type HardwareInfo struct {
	Timestamp string     `json:"timestamp"`
	System    SystemInfo `json:"system"`
	CPU       CPUInfo    `json:"cpu"`
	Disk      DiskInfo   `json:"disk"`
}

type SystemInfo struct {
	OperatingSystem          string `json:"Operating System"`
	Platform                 string `json:"Platform"`
	Hostname                 string `json:"Hostname"`
	NumberOfProcessesRunning int    `json:"Number of processes running"`
	TotalMemory              string `json:"Total memory"`
	FreeMemory               string `json:"Free memory"`
	PercentageUsedMemory     string `json:"Percentage used memory"`
}

type CPUInfo struct {
	ModelName string    `json:"Model Name"`
	Family    string    `json:"Family"`
	Speed     string    `json:"Speed"`
	Cores     []float64 `json:"cores"`
}

type DiskInfo struct {
	TotalDiskSpace           string `json:"Total disk space"`
	UsedDiskSpace            string `json:"Used disk space"`
	FreeDiskSpace            string `json:"Free disk space"`
	PercentageDiskSpaceUsage string `json:"Percentage disk space usage"`
}

func GetSystem() (SystemInfo, error) {
	runtimeOS := runtime.GOOS

	memStat, err := mem.VirtualMemory()
	if err != nil {
		return SystemInfo{}, err
	}

	hostStat, err := host.Info()
	if err != nil {
		return SystemInfo{}, err
	}

	return SystemInfo{
		OperatingSystem:          runtimeOS,
		Platform:                 hostStat.Platform,
		Hostname:                 hostStat.Hostname,
		NumberOfProcessesRunning: int(hostStat.Procs),
		TotalMemory:              fmt.Sprintf("%d", memStat.Total),
	}, nil
}

func GetCPU() (CPUInfo, error) {
	cpuStat, err := cpu.Info()
	if err != nil {
		return CPUInfo{}, err
	}

	persentage, err := cpu.Percent(0, false)
	if err != nil {
		return CPUInfo{}, err
	}

	return CPUInfo{
		ModelName: cpuStat[0].ModelName,
		Family:    cpuStat[0].Family,
		Speed:     fmt.Sprintf("%d", cpuStat[0].Mhz),
		Cores:     persentage,
	}, nil
}

func GetDisk() (DiskInfo, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return DiskInfo{}, err
	}

	return DiskInfo{
		TotalDiskSpace:           fmt.Sprintf("%d", diskStat.Total),
		UsedDiskSpace:            fmt.Sprintf("%d", diskStat.Used),
		FreeDiskSpace:            fmt.Sprintf("%d", diskStat.Free),
		PercentageDiskSpaceUsage: fmt.Sprintf("%.2f", diskStat.UsedPercent),
	}, nil
}
