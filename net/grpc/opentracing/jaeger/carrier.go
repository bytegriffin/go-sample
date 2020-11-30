package jaeger

import "google.golang.org/grpc/metadata"

/**
  自定义 Carrier
  官方propagation中只提供了两种Carrier：TextMapCarrier 和 HTTPHeadersCarrier。
  如果要实现自定义的carrier，就必须要实现 TextMapWriter 和 TextMapReader 接口
*/
type MDCarrier struct {
	metadata.MD
}

// 需要实现 opentracing.TextMapReader 接口
func (m MDCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, strs := range m.MD {
		for _, v := range strs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// 需要实现 opentracing.TextMapWriter 接口
func (m MDCarrier) Set(key, val string) {
	m.MD[key] = append(m.MD[key], val)
}
