package kitex_middleware

import (
	"context"
	"ebl_api/common/render"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/utils"
)

func ServerLogMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request, response interface{}) error {
		if arg, ok := request.(utils.KitexArgs); ok {
			if req := arg.GetFirstArgument(); req != nil {
				logger.CtxInfof(ctx, "[Get Request] Req = %v", render.Render(req))
			}
		}
		err := next(ctx, request, response)
		if result, ok := response.(utils.KitexResult); ok {
			if resp := result.GetResult(); resp != nil {
				logger.CtxInfof(ctx, "[Send Response] Resp = %v", render.Render(resp))
			}
		}
		return err
	}
}

func ClientLogMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request, response interface{}) error {
		if arg, ok := request.(utils.KitexArgs); ok {
			if req := arg.GetFirstArgument(); req != nil {
				logger.CtxInfof(ctx, "[Send Request] Req = %v", render.Render(req))
			}
		}
		err := next(ctx, request, response)
		if result, ok := response.(utils.KitexResult); ok {
			if resp := result.GetResult(); resp != nil {
				logger.CtxInfof(ctx, "[Get Response] Resp = %v", render.Render(resp))
			}
		}
		return err
	}
}
