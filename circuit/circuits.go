package http

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"time"
)

type CircuitFunc func() error
type CircuitErrorFilter func(error) (bool, error)

type CircuitConfig struct {
	Name                   string
	Timeout                int
	MaxConcurrentRequests  int
	ErrorPercentThreshold  int
	RequestVolumeThreshold int
	SleepWindow            int
	Commands               []string
}

type Circuit struct {
	config CircuitConfig
}

func NewCircuit(c CircuitConfig) *Circuit {
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

	return &Circuit{
		config: c,
	}
}

func (c *Circuit) Do(command string, fu CircuitFunc, fallback func(error) error, fi ...CircuitErrorFilter) error {
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

func (c *Circuit) DoR(command string, fu CircuitFunc, fallback func(error) error, retry int, delay time.Duration, fi ...CircuitErrorFilter) error {
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
