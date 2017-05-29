package main

import (
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/disk"
	"fmt"
)

type StatSet struct {
	load float64
	io_read float64
	io_write float64
	net float64
}

type StatsWeight struct {
	load_w float64
	io_read_w float64
	io_write_w float64
	net_w float64
}

func getAllStats() (StatSet) {
	avg, loadErr := load.Avg()
	if (loadErr != nil) {
		panic(loadErr)
	}

	ret, diskErr := disk.IOCounters()
	if (diskErr != nil) {
		panic(diskErr)
	}
	fmt.Print(ret)

	return StatSet{avg.Load1,0,0,0}
}

func (s StatSet) weight(w StatsWeight) (float64) {
	ret := float64(0)
	ret += s.load * w.load_w
	ret += s.io_read * w.io_read_w
	ret += s.io_write * w.io_write_w
	ret += s.net * w.net_w

	return ret
}