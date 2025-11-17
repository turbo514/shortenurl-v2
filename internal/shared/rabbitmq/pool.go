package rabbitmq

//type ChannelPool struct {
//	conn     *amqp.Connection
//	channels chan *amqp.Channel
//	size     int
//	closeCh  chan struct{}
//}
//
//func NewChannelPool(conn *amqp.Connection, poolSize int) (*ChannelPool, error) {
//	if poolSize <= 0 {
//		return nil, errors.New("poolSize 必须大于0")
//	}
//
//	pool := &ChannelPool{
//		conn:     conn,
//		channels: make(chan *amqp.Channel, poolSize),
//		closeCh:  make(chan struct{}),
//		size:     poolSize,
//	}
//
//	for range poolSize {
//		channel, err := conn.Channel()
//		if err != nil {
//			for {
//				ch, ok := <-pool.channels
//				if !ok {
//					break
//				} else {
//					ch.Close()
//				}
//			}
//			return nil, fmt.Errorf("创建Channel Pool失败: %w", err)
//		}
//		pool.channels <- channel
//	}
//
//	return pool, nil
//}
//
//func (p *ChannelPool) Get(ctx context.Context) (*amqp.Channel, error) {
//	select {
//	case <-p.closeCh:
//		return nil, fmt.Errorf("ChannelPool已关闭")
//	case <-ctx.Done():
//		return nil, fmt.Errorf("获取channel超时: %w", ctx.Err())
//	case channel := <-p.channels:
//		return channel, nil
//	}
//}
//
//func (p *ChannelPool) Put(ctx context.Context, channel *amqp.Channel) {
//	select {
//	case <-p.closeCh:
//		channel.Close()
//		return
//	case <-ctx.Done():
//		channel.Close()
//		return
//	case p.channels <- channel:
//		return
//	default:
//		// 连接池已满
//		channel.Close()
//	}
//}
//
//func (p *ChannelPool) Close() {
//	close(p.closeCh)
//}
