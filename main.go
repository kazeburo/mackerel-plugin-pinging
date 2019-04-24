package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"runtime"
	"sort"
	"time"

	flags "github.com/jessevdk/go-flags"
	fping "github.com/tatsushid/go-fastping"
)

// Version by Makefile
var Version string

type cmdOpts struct {
	Host      string `long:"host" description:"Hostname to ping" required:"true"`
	MaxRtt    int    `long:"max-rtt" default:"200" description:"wait time. Max RTT(ms)"`
	Count     int    `long:"count" default:"10" description:"Count Sending ping"`
	KeyPrefix string `long:"key-prefix" description:"Metric key prefix" required:"true"`
	Version   bool   `short:"v" long:"version" description:"Show version"`
}

func round(f float64) int64 {
	return int64(math.Round(f)) - 1
}

func getStats(opts cmdOpts) error {
	ra, err := net.ResolveIPAddr("ip", opts.Host)
	if err != nil {
		return err
	}

	var rtts sort.Float64Slice
	var t float64
	s := float64(0)
	e := float64(0)
	pinger := fping.NewPinger()
	pinger.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		rttMilliSec := float64(rtt.Nanoseconds()) / 1000.0 / 1000.0
		rtts = append(rtts, rttMilliSec)
		t += rttMilliSec
		s++
	}
	pinger.AddIPAddr(ra)
	pinger.MaxRTT = time.Millisecond * time.Duration(opts.MaxRtt)
	for i := 0; i < opts.Count; i++ {
		err := pinger.Run()
		if err != nil {
			log.Printf("%v", err)
			e++
		}
	}
	now := uint64(time.Now().Unix())
	fmt.Printf("pinging.%s_rtt_count.success\t%f\t%d\n", opts.KeyPrefix, s, now)
	fmt.Printf("pinging.%s_rtt_count.error\t%f\t%d\n", opts.KeyPrefix, e, now)
	if s > 0 {
		fmt.Printf("pinging.%s_rtt_ms.average\t%f\t%d\n", opts.KeyPrefix, t/s, now)
		fmt.Printf("pinging.%s_rtt_ms.max\t%f\t%d\n", opts.KeyPrefix, rtts[round(s)], now)
		fmt.Printf("pinging.%s_rtt_ms.90_percentile\t%f\t%d\n", opts.KeyPrefix, rtts[round(s*0.90)], now)
	}
	return nil
}

func main() {
	os.Exit(_main())
}

func _main() int {
	opts := cmdOpts{}
	psr := flags.NewParser(&opts, flags.Default)
	_, err := psr.Parse()
	if err != nil {
		return 1
	}
	if opts.Version {
		fmt.Printf(`%s %s
Compiler: %s %s
`,
			os.Args[0],
			Version,
			runtime.Compiler,
			runtime.Version())
		return 0
	}

	err = getStats(opts)
	if err != nil {
		log.Printf("%v", err)
		return 1
	}
	return 0
}
