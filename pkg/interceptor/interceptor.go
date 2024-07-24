package interceptor

import (
	"context"

	"github.com/digisata/invitation-service/pkg/constans"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type InterceptorManager interface {
	Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)
	ClientRequestLoggerInterceptor() func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error
}

type interceptorManager struct {
	logger *zap.SugaredLogger
}

func NewInterceptorManager(logger *zap.SugaredLogger) *interceptorManager {
	return &interceptorManager{
		logger: logger,
	}
}

func (im interceptorManager) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	reply, err := handler(ctx, req)
	if err != nil {
		im.logger.Errorw(constans.ERROR,
			"method", info.FullMethod,
			"request", req,
			"error", err.Error(),
		)

		return reply, err
	}

	im.logger.Infow(constans.INFO,
		"method", info.FullMethod,
		"request", req,
		"error", nil,
	)

	return reply, err

}

func (im interceptorManager) ClientRequestLoggerInterceptor() func(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		im.logger.Infow(constans.INFO,
			"method", method,
			"request", req,
			"error", nil,
		)

		return err
	}
}
