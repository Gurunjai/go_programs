package main

import (
	"testing"
	"fmt"
)

func TestGetMem(t *testing.T) {
	fmt.Printf("Read Mem Stat is: %v\n", getMemStats())
}

func TestGetNetStats(t *testing.T) {
	fmt.Printf("Read Network Stat is: %v\n", getNetStats())
}