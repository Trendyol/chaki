package client

import (
	"context"
	"github.com/Trendyol/chaki/config"
	"github.com/avast/retry-go"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
)

type Base struct {
	name   string
	driver *resty.Client
	cb     *cbCfg
	retry  *retryCfg
}

type Factory struct {
	cfg          *config.Config
	baseWrappers []DriverWrapper
}

func NewFactory(cfg *config.Config, wrappers []DriverWrapper) *Factory {
	return &Factory{
		cfg:          cfg,
		baseWrappers: wrappers,
	}
}

func (f *Factory) Get(name string, errDecoder ErrDecoder, driverWrappers ...DriverWrapper) *Base {
	clientCfg := f.cfg.Of("client").Of(name)
	cb := cbFromCfg(clientCfg.Of("circuit-breaker"))
	r := retryFromCfg(clientCfg.Of("retry"))
	return &Base{
		driver: newDriverBuilder(clientCfg).
			AddErrDecoder(errDecoder).
			AddUpdaters(f.baseWrappers...).
			AddUpdaters(driverWrappers...).
			build(),
		name:  name,
		cb:    cb,
		retry: r,
	}
}

func (r *Base) Request(ctx context.Context) *resty.Request {
	return r.driver.R().SetContext(ctx)
}

func (r *Base) Send(request *resty.Request, uri string, method string) (*resty.Response, error) {
	if r.cb.Enabled {
		resp, err := r.cb.Cb.Execute(func() (interface{}, error) {
			if r.retry.Enabled {
				return r.sendWithRetry(request, uri, method)
			}
			return r.send(request, uri, method)
		})
		if err != nil {
			return nil, err
		}
		return resp.(*resty.Response), err
	}

	if r.retry.Enabled {
		return r.sendWithRetry(request, uri, method)
	}
	return r.send(request, uri, method)
}

func (r *Base) sendWithRetry(request *resty.Request, uri string, method string) (*resty.Response, error) {
	var resp *resty.Response
	var err error
	err = retry.Do(
		func() error {
			resp, err = r.send(request, uri, method)
			return err
		},
		retry.Attempts(r.retry.Attempts),
		retry.Delay(r.retry.Delay),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Retry #%d for request to %s failed: %s", n+1, uri, err)
		}),
	)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (r *Base) send(request *resty.Request, uri string, method string) (*resty.Response, error) {
	switch method {
	case http.MethodGet:
		return request.Get(uri)
	case http.MethodPost:
		return request.Post(uri)
	case http.MethodPut:
		return request.Put(uri)
	case http.MethodDelete:
		return request.Delete(uri)
	case http.MethodPatch:
		return request.Patch(uri)
	case http.MethodHead:
		return request.Head(uri)
	case http.MethodOptions:
		return request.Options(uri)
	default:
		return nil, ErrUnsupportedMethod
	}
}
