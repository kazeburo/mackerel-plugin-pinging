package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	ping "github.com/digineo/go-ping"
	flags "github.com/jessevdk/go-flags"
)

// Version by Makefile
var Version string

type cmdOpts struct {
	Host      string `long:"host" description:"Hostname to ping" required:"true"`
	Timeout   int    `long:"timeout" default:"1000" description:"timeout millisec per ping"`
	Interval  int    `long:"interval" default:"10" description:"sleep millisec after every ping"`
	Count     int    `long:"count" default:"10" description:"Count Sending ping"`
	KeyPrefix string `long:"key-prefix" description:"Metric key prefix" required:"true"`
	Version   bool   `short:"v" long:"version" description:"Show version"`
}

func round(f float64) int64 {
	return int64(math.Round(f)) - 1
}

func getStats(opts cmdOpts) error {
	var ra *net.IPAddr
	var pinger *ping.Pinger
	if strings.Index(opts.Host, ":") != -1 {
		r, err := net.ResolveIPAddr("ip6", opts.Host)
		if err != nil {
			return err
		}
		ra = r
		p, err := ping.New("", "::")
		if err != nil {
			return err
		}
		pinger = p
	} else {
		r, err := net.ResolveIPAddr("ip4", opts.Host)
		if err != nil {
			return err
		}
		ra = r
		p, err := ping.New("0.0.0.0", "")
		if err != nil {
			return err
		}
		pinger = p
	}

	defer pinger.Close()

	var rtts sort.Float64Slice
	var t float64
	s := float64(0)
	e := float64(0)

	// preflight
	_, err := pinger.Ping(ra, time.Millisecond * time.Duration(opts.Timeout))
	if err != nil {
		log.Printf("error in preflight: %v", err)
	}

	for i := 0; i < opts.Count; i++ {
		time.Sleep(time.Millisecond * time.Duration(opts.Interval))
		rtt, err := pinger.Ping(ra, time.Millisecond * time.Duration(opts.Timeout))
		if err != nil {
			log.Printf("%v", err)
			e++
			continue
		}
		rttMilliSec := float64(rtt.Nanoseconds()) / 1000.0 / 1000.0
		rtts = append(rtts, rttMilliSec)
		t += rttMilliSec
		s++
	}
	sort.Sort(rtts)
	now := uint64(time.Now().Unix())
	fmt.Printf("pinging.%s_rtt_count.success\t%f\t%d\n", opts.KeyPrefix, s, now)
	fmt.Printf("pinging.%s_rtt_count.error\t%f\t%d\n", opts.KeyPrefix, e, now)
	if s > 0 {
		fmt.Printf("pinging.%s_rtt_ms.max\t%f\t%d\n", opts.KeyPrefix, rtts[round(s)], now)
		fmt.Printf("pinging.%s_rtt_ms.min\t%f\t%d\n", opts.KeyPrefix, rtts[0], now)
		fmt.Printf("pinging.%s_rtt_ms.average\t%f\t%d\n", opts.KeyPrefix, t/s, now)
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
