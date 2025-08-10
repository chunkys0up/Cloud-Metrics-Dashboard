package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

// get the timed cpu usage every second
func sampleCPU() {
	for {
		percent, err := cpu.Percent(time.Second, false)
		if err == nil && len(percent) > 0 {
			mu.Lock()
			cpu_usage = percent[0]
			mu.Unlock()
		}
	}
}

func sampleMemory() float64 {
	v, _ := mem.VirtualMemory()

	return v.UsedPercent
}

func sampleDisk() float64 {
	v, _ := disk.Usage("/")

	return v.UsedPercent
}

func sampleBytes() {
	kilobytes_per_second := float64(1) / 1000

	for {
		initialIoData, err := net.IOCounters(false)
		if err != nil {
			fmt.Println("Error getting intial I/O data:", err)
			return
		}

		initialTime := time.Now()

		time.Sleep(1 * time.Second)

		endIoData, err := net.IOCounters(false)
		if err != nil {
			fmt.Println("Error getting end I/O data:", err)
			return
		}

		endTime := time.Now()

		duration := endTime.Sub(initialTime).Seconds()
		bytesRecv := float64(endIoData[0].BytesRecv - initialIoData[0].BytesRecv)
		bytesSent := float64(endIoData[0].BytesSent - initialIoData[0].BytesSent)

		mu.Lock()
		bytesRecvRate = bytesRecv / duration * float64(kilobytes_per_second)
		bytesSentRate = bytesSent / duration  * float64(kilobytes_per_second)
		mu.Unlock()
	}
}

