package client

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"time"
)

type circuitFunc func() error
type CircuitErrorFunc func(error) error
type CircuitErrorFilter func(error) (bool, error)

type circuitConfig struct {
	Name                   string
	Timeout                int
	MaxConcurrentRequests  int
	ErrorPercentThreshold  int
	RequestVolumeThreshold int
	SleepWindow            int
	Commands               []string
}

type circuit struct {
	config *circuitConfig
}

func newCircuit(c *circuitConfig) *circuit {
	hystrixConfig := hystrix.CommandConfig{
		Timeout:                c.Timeout,
		MaxConcurrentRequests:  c.MaxConcurrentRequests,
		ErrorPercentThreshold:  c.ErrorPercentThreshold,
		RequestVolumeThreshold: c.RequestVolumeThreshold,
		SleepWindow:            c.SleepWindow,
	}

	for _, command := range c.Commands {
		hystrix.ConfigureCommand(fmt.Sprintf("%s:%s", c.Name, command), hystrixConfig)
	}

	return &circuit{
		config: c,
	}
}

func (c *circuit) do(command string, fu circuitFunc, fallback func(error) error, fi ...CircuitErrorFilter) error {
	var e error
	var ok bool

	function := func() error {

		err := fu()

		for _, filter := range fi {
			if ok, e = filter(err); ok {
				return err
			}
		}

		if len(fi) > 0 {
			return nil
		}

		return err
	}

	hystrixErr := hystrix.Do(fmt.Sprintf("%s:%s", c.config.Name, command), function, fallback)

	if hystrixErr != nil {
		return hystrixErr
	}

	if e != nil {
		return e
	}

	return nil
}

func (c *circuit) doR(command string, fu circuitFunc, fallback func(error) error, retry int, delay time.Duration, fi ...CircuitErrorFilter) error {
	var e error
	var ok bool

	filter := func() error {
		err := fu()

		for _, filter := range fi {
			if ok, e = filter(err); ok {
				return err
			}
		}

		if len(fi) > 0 {
			return nil
		}

		return err
	}

	function := func() error {
		err := filter()

		for c := 0; c < retry && err != nil; c++ {
			time.Sleep(delay)
			err = filter()
		}

		return err
	}

	hystrixErr := hystrix.Do(fmt.Sprintf("%s:%s", c.config.Name, command), function, fallback)

	if hystrixErr != nil {
		return hystrixErr
	}

	if e != nil {
		return e
	}

	return nil
}

func defaultCircuitErrorFunc(err error) error {
	return err
}
