package hardware

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

func GetSystem() (string, error) {
	runtimeOS := runtime.GOOS

	memStat, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}

	hostStat, err := host.Info()
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf(
		"OS: %s\nHost Name: %s\nTotal Memory: %d\nUsed Memory: %d\n",
		runtimeOS, hostStat.Hostname, memStat.Total, memStat.Used)

	return output, nil
}

func GetCPU() (string, error) {
	cpuStat, err := cpu.Info()
	if err != nil {
		return "", err
	}

	persentage, err := cpu.Percent(0, false)
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("CPU: %s\nCPU Usage: %.2f%%\n", cpuStat[0].ModelName, persentage[0])

	return output, nil
}

func GetDisk() (string, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("Total Disk Space: %d\nDisk Usage: %.2fGB\n Free Disk Space: %d", diskStat.Total, diskStat.UsedPercent, diskStat.Free)
	return output, nil
}
