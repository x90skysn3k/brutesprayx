package modules

import (
	"fmt"
	"time"

	"github.com/GoFeGroup/gordp"
	"github.com/GoFeGroup/gordp/proto/bitmap"
)

type RDPProcessor struct{}

func (p *RDPProcessor) ProcessBitmap(opt *bitmap.Option, bm *bitmap.BitMap) {
	// Add your custom logic here to process the bitmap
}

type Processor interface {
	ProcessBitmap(*bitmap.Option, *bitmap.BitMap)
}

func BruteRDP(host string, port int, user, password string) bool {
	const defaultTimeout = 5 * time.Second
	timer := time.NewTimer(defaultTimeout)
	defer timer.Stop()
	type result struct {
		client *gordp.Client
		err    error
	}
	done := make(chan result)

	option := &gordp.Option{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		UserName: user,
		Password: password,
	}
	client := gordp.NewClient(option)
	go func() {
		processor := &RDPProcessor{}
		err := client.Run(processor)
		done <- result{client: client, err: err}
	}()

	select {
	case <-timer.C:
		return false
	case result := <-done:
		if result.err != nil {
			return false
		}
		return true
	}
}
