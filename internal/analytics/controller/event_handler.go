package controller

import (
	"context"
	"github.com/turbo514/shortenurl-v2/analytics/cqrs/command"
	"github.com/turbo514/shortenurl-v2/shared/mylog"
	mytrace "github.com/turbo514/shortenurl-v2/shared/trace"
)

// TODO: 实现泛型

//	type ClickEventHandler struct {
//		eventChannel <-chan ClickEventWithAcker
//		buffSize     int
//		handler      *command.CreateClickEventHandler
//		closeCh      chan struct{}
//	}
//
//	func NewClickEventHandler(eventChannel <-chan ClickEventWithAcker, buffSize int, handler *command.CreateClickEventHandler) *ClickEventHandler {
//		h := &ClickEventHandler{
//			eventChannel: eventChannel,
//			buffSize:     buffSize,
//			handler:      handler,
//			closeCh:      make(chan struct{}),
//		}
//
//		return h
//	}
//
//	func (h *ClickEventHandler) Close() {
//		close(h.closeCh)
//	}
//
//	func (h *ClickEventHandler) Start() error {
//		buffer := make([]*domain.ClickEvent, 0, h.buffSize)
//		var last IAcker
//		for {
//			timeout := time.After(time.Second)
//			select {
//			case <-h.closeCh:
//				return nil
//			case event := <-h.eventChannel:
//				mylog.GetLogger().Debug("获取到点击事件", "event", event)
//				buffer = append(buffer, event.Event)
//				last = event.Acker
//				if len(buffer) >= h.buffSize {
//					//tr := mytrace.GetTracer()
//					//// 使用 context.Background() 表示没有父 trace
//					//ctx := context.Background()
//					//// 开启新的 trace（根 span）
//					//ctx, span := tr.Start(ctx, "批量写入短链接点击事件")
//					ctx := context.Background()
//					if err := h.handler.Handle(ctx, command.CreateClickEventCommand{Events: buffer}); err != nil {
//						last.Nack()
//					} else {
//						last.Ack()
//					}
//					buffer = buffer[:0]
//					last = nil
//					//span.End()
//				}
//			case <-timeout:
//				if len(buffer) > 0 {
//					//tr := mytrace.GetTracer()
//					//// 使用 context.Background() 表示没有父 trace
//					//ctx := context.Background()
//					//// 开启新的 trace（根 span）
//					//ctx, span := tr.Start(ctx, "批量写入短链接点击事件")
//					ctx := context.Background()
//					if err := h.handler.Handle(ctx, command.CreateClickEventCommand{Events: buffer}); err != nil {
//						last.Nack()
//					} else {
//						last.Ack()
//					}
//					//span.End()
//					buffer = buffer[:0]
//					last = nil
//				}
//			}
//		}
//	}
type ClickEventHandler struct {
	Recs    chan ClickEventWithAcker
	quitCh  chan struct{}
	done    chan struct{}
	handler *command.CreateClickEventHandler
}

func NewClickEventHandler(handler *command.CreateClickEventHandler) *ClickEventHandler {
	return &ClickEventHandler{
		Recs:    make(chan ClickEventWithAcker),
		quitCh:  make(chan struct{}),
		done:    make(chan struct{}),
		handler: handler,
	}
}

func (h *ClickEventHandler) Consume() {
	logger := mylog.GetLogger()

	defer close(h.done) // 消费完成关闭 done channel

	for {
		select {
		case <-h.quitCh:
			return
		case rec := <-h.Recs:
			tr := mytrace.GetTracer()
			ctx := context.Background()
			ctx, span := tr.Start(ctx, "ClickEventHandler.Consume")

			for i := range rec.Event {
				logger.Debug("监听到点击事件", "event", rec.Event[i])
			}

			if err := h.handler.Handle(ctx, command.CreateClickEventCommand{Events: rec.Event}); err != nil {
				logger.Warn("事件处理失败")
				rec.Acker.Nack(context.Background())
			} else {
				rec.Acker.Ack(context.Background())
			}
			span.End()
		}
	}

}

func (h *ClickEventHandler) Quit() {
	close(h.quitCh)
}

func (h *ClickEventHandler) Done() <-chan struct{} {
	return h.done
}
