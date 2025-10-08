package channels

import (
	"context"
	"inpayos/internal/protocol"
	"inpayos/internal/task"
)

func init() {
	// 注册渠道同步任务
	task.RegisterHandler("channel_sync", syncHandler)
}

// syncHandler 渠道同步处理器
func syncHandler(ctx context.Context, params protocol.MapData) error {
	LoadChannelOpenApiService()
	return nil
}
