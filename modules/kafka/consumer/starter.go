package consumer

type Starter interface {
	Start() error
	Stop() error
}

//func buildConsumeFn(cf ConsumeFn, interceptors []Interceptor) ConsumeFn {
//	if len(interceptors) == 0 {
//		return cf
//	}
//
//	next := buildNextFunc(cf, interceptors, 1)
//	return func(m *kafkalib.Message) error {
//		return interceptors[0].Intercept(m, next)
//	}
//
//}
//
//func buildNextFunc(cf ConsumeFn, interceptors []Interceptor, i int) ConsumeFn {
//	if i > len(interceptors)-1 {
//		return func(m *kafkalib.Message) error { return cf(m) }
//	}
//
//	next := buildNextFunc(cf, interceptors, i+1)
//	return func(m *kafkalib.Message) error {
//		return interceptors[i].Intercept(m, next)
//	}
//}
//
//func buildBatchConsumeFn(cf BatchConsumeFn, interceptors []BatchConsumerInterceptor) BatchConsumeFn {
//	if len(interceptors) == 0 {
//		return cf
//	}
//
//	next := buildBatchConsumeNextFunc(cf, interceptors, 1)
//	return func(m []*kafkalib.Message) error {
//		return interceptors[0].Intercept(m, next)
//	}
//
//}
//
//func buildBatchConsumeNextFunc(cf BatchConsumeFn, interceptors []BatchConsumerInterceptor, i int) BatchConsumeFn {
//	if i > len(interceptors)-1 {
//		return func(m []*kafkalib.Message) error { return cf(m) }
//	}
//
//	next := buildBatchConsumeNextFunc(cf, interceptors, i+1)
//	return func(m []*kafkalib.Message) error {
//		return interceptors[i].Intercept(m, next)
//	}
//}
