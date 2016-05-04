package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
)

// Ipbeat represent a beat to send IP to ELKstack
type Ipbeat struct {
	period time.Duration

	IpbeatConfig ConfigSettings
	events       publisher.Client

	done chan struct{}
}

// New creates a new Ipbeat
func New() *Ipbeat {
	return &Ipbeat{}
}

// Config reads configuration for Ipbeat from ipbeat.yml
func (ipbeat *Ipbeat) Config(beat *beat.Beat) error {
	err := beat.RawConfig.Unpack(&ipbeat.IpbeatConfig)

	if err != nil {
		logp.Err("IPbeat failed to read configuration file: %v", err)

		return err
	}

	if ipbeat.IpbeatConfig.Input != nil {
		return fmt.Errorf("'input' will soon be deprecated. Use 'ipbeat' instead")
	}

	ipbeatConfig := ipbeat.IpbeatConfig.Ipbeat

	if ipbeatConfig.Period != nil {
		ipbeat.period = time.Duration(*ipbeatConfig.Period) * time.Second
	} else {
		ipbeat.period = 10 * time.Second
	}

	logp.Debug("ipbeat", "IPBeat has been configured")
	logp.Debug("ipbeat", "Period: %v", ipbeat.period)

	return nil
}

// Setup makes ipbeat ready to run
func (ipbeat *Ipbeat) Setup(beat *beat.Beat) error {
	ipbeat.events = beat.Publisher.Connect()
	ipbeat.done = make(chan struct{})

	return nil
}

// Run is the main function that sends IP data to ELK
func (ipbeat *Ipbeat) Run(beat *beat.Beat) error {
	ticker := time.NewTicker(ipbeat.period)

	defer ticker.Stop()

	var ip IP

	for {
		select {
		case <-ipbeat.done:
			return nil
		case <-ticker.C:

			timerStart := time.Now()

			gotIP, err := ip.RetrieveIP()

			if err != nil {
				logp.Err("Error while getting IP: %v", err)
			}

			fmt.Printf("%v", gotIP)

			event := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"ip":         gotIP,
			}

			ipbeat.events.PublishEvent(event)

			timerEnd := time.Now()

			if timerEnd.Sub(timerStart).Nanoseconds() > ipbeat.period.Nanoseconds() {
				logp.Warn("Processing took more than one period.")
			}
		}
	}
}

// Cleanup cleans memory & stuff before quitting
func (ipbeat *Ipbeat) Cleanup(beat *beat.Beat) error {
	return nil
}

// Stop stops the ipbeat
func (ipbeat *Ipbeat) Stop() {
	logp.Info("Stopping IPBeat")
	close(ipbeat.done)
	ipbeat.events.Close()
}
