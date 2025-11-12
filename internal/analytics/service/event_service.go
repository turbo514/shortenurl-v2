package service

//type EventService struct {
//	createClickEventHandler *command.CreateClickEventHandler
//}
//
//func NewEventService(createClickHandler *command.CreateClickEventHandler, clickCounter domain.IClickCounterWriter) *EventService {
//	return &EventService{
//		createClickEventHandler: createClickHandler,
//	}
//}
//
//func (s *EventService) HandleClickEvents(ctx context.Context, events []*domain.ClickEvent) error {
//	if err := s.createClickEventHandler.Handle(ctx, command.CreateClickEventCommand{Events: events}); err != nil {
//		return fmt.Errorf("写入点击事件失败: %w", err)
//	}
//
//	fmt.Printf("%+v\n", events)
//
//	return nil
//}
