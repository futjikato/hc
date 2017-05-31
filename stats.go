package main

import (
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/net"
	"time"
)

type StatSet struct {
	Ts time.Time `json:"time"`

	Load         float64          `json:"load"`
	IoRead       map[string]int64 `json:"io_read"`
	IoReadTotal  int64            `json:"io_read_total"`
	IoWrite      map[string]int64 `json:"io_write"`
	IoWriteTotal int64            `json:"io_write_total"`
	NetSent      map[string]int64 `json:"net_sent"`
	NetSentTotal int64            `json:"net_sent_total"`
	NetRecv      map[string]int64 `json:"net_recv"`
	NetRecvTotal int64            `json:"net_recv_total"`
}

type StatConfig struct {
	load bool
	io   bool
	net  bool

	lastIoReadCount  map[string]int64
	lastIoWriteCount map[string]int64
	lastNetSentBytes map[string]int64
	lastNetRecvBytes map[string]int64
}

func getStats(c *StatConfig) *StatSet {
	s := &StatSet{}
	s.Ts = time.Now()

	if c.load {
		s.Load = getLoad()
	}

	if c.io {
		ret, err := disk.IOCounters()
		if err != nil {
			panic(err)
		}

		s.IoRead = getReadCounts(ret, c)
		for _, rv := range s.IoRead {
			s.IoReadTotal += rv
		}
		s.IoWrite = getWriteCounts(ret, c)
		for _, wv := range s.IoWrite {
			s.IoWriteTotal += wv
		}
	}

	if c.net {
		netRet, netErr := net.IOCounters(false)
		if netErr != nil {
			panic(netErr)
		}

		s.NetSent = getSentBytes(netRet, c)
		for _, sv := range s.NetSent {
			s.NetSentTotal += sv
		}

		s.NetRecv = getRecvBytes(netRet, c)
		for _, rv := range s.NetRecv {
			s.NetRecvTotal += rv
		}
	}

	return s
}

func getLoad() float64 {
	avg, loadErr := load.Avg()
	if loadErr != nil {
		panic(loadErr)
	}

	return avg.Load1
}

func getReadCounts(ret map[string]disk.IOCountersStat, c *StatConfig) map[string]int64 {
	r := make(map[string]int64)

	for key, d := range ret {
		if c.lastIoReadCount[key] > 0 {
			r[key] = int64(d.ReadCount) - c.lastIoReadCount[key]
		}
		c.lastIoReadCount[key] = int64(d.ReadCount)
	}

	return r
}

func getWriteCounts(ret map[string]disk.IOCountersStat, c *StatConfig) map[string]int64 {
	r := make(map[string]int64)

	for key, d := range ret {
		if c.lastIoWriteCount[key] > 0 {
			r[key] = int64(d.WriteCount) - c.lastIoWriteCount[key]
		}
		c.lastIoWriteCount[key] = int64(d.WriteCount)
	}

	return r
}

func getSentBytes(ret []net.IOCountersStat, c *StatConfig) map[string]int64 {
	r := make(map[string]int64)

	for _, d := range ret {
		if c.lastNetSentBytes[d.Name] > 0 {
			r[d.Name] = int64(d.BytesSent) - c.lastNetSentBytes[d.Name]
		}
		c.lastNetSentBytes[d.Name] = int64(d.BytesSent)
	}

	return r
}

func getRecvBytes(ret []net.IOCountersStat, c *StatConfig) map[string]int64 {
	r := make(map[string]int64)

	for _, d := range ret {
		if c.lastNetRecvBytes[d.Name] > 0 {
			r[d.Name] = int64(d.BytesRecv) - c.lastNetRecvBytes[d.Name]
		}
		c.lastNetRecvBytes[d.Name] = int64(d.BytesRecv)
	}

	return r
}
