package tracing

// import (
// 	"io"
// 	"time"

// 	"github.com/opentracing/opentracing-go"
// 	"github.com/uber/jaeger-client-go/config"
// )

// type OpenTracer struct {
// 	ServiceName string
// 	Address     string
// 	Tracer      opentracing.Tracer
// }

// func NewOpenTracer(serviceName, address string) (opentracing.Tracer, io.Closer, error) {
// 	cfg := config.Configuration{
// 		Sampler: &config.SamplerConfig{
// 			Type:  "const",
// 			Param: 1,
// 		},
// 		Reporter: &config.ReporterConfig{
// 			LogSpans:            false,
// 			BufferFlushInterval: 1 * time.Second,
// 			LocalAgentHostPort:  address,
// 		},
// 	}
// 	return cfg.New(
// 		"ExampleTracingMiddleware", //service name
// 	)
// }
