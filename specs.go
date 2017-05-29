package main

import "fmt"

type Spec interface {
	getStats()
}

type CpuSpec struct {}
func (s CpuSpec) getStats() {
	stats := getAllStats()
	health := stats.weight(StatsWeight{1,0.2,0.2,0.6})

	fmt.Printf("Health is %f.4", health)
}

type NetSpec struct {}

func (s NetSpec) getStats() {
	stats := getAllStats()
	health := stats.weight(StatsWeight{0.5,0.2,0.2,1})

	fmt.Printf("Health is %f.4", health)
}

type IoSpec struct {}

func (s IoSpec) getStats() {
	stats := getAllStats()
	health := stats.weight(StatsWeight{0.6,1,1,0.2})

	fmt.Printf("Health is %f.4", health)
}