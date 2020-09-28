package main

import (
	"fmt"
	"runtime"
	"math/rand"
	"sync"
	"time"
)

const MAX_SAMPLES = 4

type Metric interface {
	Dump()
}

type valueType  struct {
	val [MAX_SAMPLES]float64
	activeIdx uint8
}

func (v *valueType) set(f float64) {
	v.val[v.activeIdx] = f
	v.activeIdx = ((v.activeIdx + 1) % MAX_SAMPLES)
}

func (v *valueType) read(idx int) float64 {
	if idx == -1 {
		i := v.activeIdx
		if i != 0 {
			i -= 1
		}
		return v.val[i]
	} else {
		return v.val[idx]
	}
}

type Mem struct {
	valueType
	f func() float64
}

type Net struct {
	Rxb valueType
	Txb valueType
	Rxr valueType
	Txr valueType
	f func() []float64
}

func StartMem(f func() float64) Metric {
	m := &Mem{ f: f }
	m.Do()
	return m
}

func (m *Mem) Do() {
	wait := make(chan bool)
	go func(m *Mem, wait chan <- bool) {
		wait <- true
		close(wait)
		for {
			select {
			case <- time.After(1 * time.Second):
				m.set(m.f())
			}
		}

	} (m, wait)

	<-wait
}

func (m *Mem) Dump() {
	fmt.Printf("Memory Values are:\t")
	for _, v := range m.val {
		fmt.Printf("%v MB, ", v)
	}
	fmt.Printf("\n")
}

func StartNet(f func() []float64) Metric {
	n := &Net{ f: f }
	n.Do()
	return n
}

func (n *Net) Do() {
	wait := make(chan bool)
	go func(n *Net, wait chan <- bool) {
		wait <- true
		close(wait)
		for {
			select {
			case <- time.After(1 * time.Second):
				out := n.f()
				if len(out) != 2 {
					panic("invalid length of response")
				}
				rx, tx := out[0], out[1]

				prx, ptx := n.Rxb.read(-1), n.Txb.read(-1)
				
				n.Rxr.set((rx - prx) * 8)
				n.Txr.set((tx - ptx) * 8)
				n.Rxb.set(rx)
				n.Txb.set(tx)
			}
		}

	} (n, wait)

	<-wait
}

func (n *Net) Dump() {
	fmt.Printf("Network Interface values are:")
	fmt.Printf("\n\tRx Bytes:\t")
	for i := 0; i < MAX_SAMPLES; i++ {
		fmt.Printf("%0.2f, ", n.Rxb.read(i))
	}

	fmt.Printf("\n\tTx Bytes: \t")
	for i := 0; i < MAX_SAMPLES; i++ {
		fmt.Printf("%0.2f, ", n.Txb.read(i))
	}

	fmt.Printf("\n\tRx Rate (bps): \t")
	for i := 0; i < MAX_SAMPLES; i++ {
		fmt.Printf("%0.2f, ", n.Rxr.read(i))
	}

	fmt.Printf("\n\tTx Rate (bps): \t")
	for i := 0; i < MAX_SAMPLES; i++ {
		fmt.Printf("%0.2f, ", n.Txr.read(i))
	}

	fmt.Printf("\n")
}

func getMemStats() float64 {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return float64(mem.Sys / (1 << 20))
}

func getNetStats() []float64 {
	// Mimic what /proc/net/dev provides
	r := rand.New(rand.NewSource(rand.Int63n((1 << 19))))
	rxB := r.NormFloat64()
	txB := r.NormFloat64()

	return []float64{rxB, txB}	
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	m := StartMem(getMemStats)
	n := StartNet(getNetStats)
	go func(m, n *Metric, wg *sync.WaitGroup) {
		defer wg.Done()

		t := time.NewTimer(40 * time.Second)

		lLoop:
		for {
			select {
			case <- t.C:
				fmt.Printf("Ended!!!!\n")
				break lLoop
			case <- time.After(10 * time.Second):
				(*m).Dump()
				(*n).Dump()
				fmt.Printf("\n")
			}
		}

	} (&m, &n, &wg)

	wg.Wait()
}