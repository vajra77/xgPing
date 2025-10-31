package probe

import (
	probing "github.com/prometheus-community/pro-bing"
)

type Peer struct {
	name      string
	code      string
	v4Address string
	v6Address string
	samples   [10]*Sample
	lastIndex int
}

func NewPeer(name, code, v4Address, v6Address string) *Peer {
	return &Peer{name, code, v4Address, v6Address, [10]*Sample{}, 0}
}

func (p *Peer) Name() string {
	return p.name
}

func (p *Peer) Code() string {
	return p.code
}

func (p *Peer) V4Address() string {
	return p.v4Address
}

func (p *Peer) V6Address() string {
	return p.v6Address
}

func (p *Peer) Ping(count int) {
	pinger, err := probing.NewPinger(p.v4Address)
	if err != nil {
		panic(err)
	}
	pinger.Count = count
	pinger.OnFinish = func(stats *probing.Statistics) {
		sample := NewSample(
			1000*stats.MinRtt.Seconds(),
			1000*stats.MaxRtt.Seconds(),
			1000*stats.AvgRtt.Seconds(),
			1000*stats.StdDevRtt.Seconds(),
			100*stats.PacketLoss,
		)
		if p.lastIndex >= 10 {
			p.lastIndex = 0
		} else {
			p.lastIndex++
		}
		p.samples[p.lastIndex] = sample
	}
}

func (p *Peer) LastSample() *Sample {
	return p.samples[p.lastIndex]
}
