package app

import (
	"context"

	"github.com/actuagent/cclaw/pkg/bus"
	"github.com/actuagent/cclaw/pkg/config"
	"github.com/actuagent/cclaw/pkg/heartbeat"
	einomodel "github.com/cloudwego/eino/components/model"
)

func StartHeartbeatService(
	ctx context.Context,
	heartbeatCfg config.HeartbeatConfig,
	chatModel einomodel.ChatModel,
	messageBus *bus.MessageBus,
) *heartbeat.HeartbeatService {
	if !heartbeatCfg.IsEnabled() {
		return nil
	}
	heartbeatService := heartbeat.NewHeartbeatService(
		heartbeatCfg.Path,
		chatModel,
		func(ctx context.Context, tasks string) error {
			logApp.Info("Heartbeat triggered", "tasks", tasks)
			messageBus.PublishInbound(ctx, &bus.InboundMessage{
				Channel:  "heartbeat",
				ChatID:   "system",
				Content:  tasks,
				Metadata: map[string]any{"type": "heartbeat"},
			})
			return nil
		},
		heartbeatCfg.Interval.Duration,
	)
	heartbeatService.Start(ctx)
	return heartbeatService
}
