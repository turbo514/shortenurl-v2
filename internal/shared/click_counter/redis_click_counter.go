package click_counter

//var _ IClickCounterReader = (*RedisClickCounterReader)(nil)
//
//type RedisClickCounterReader struct {
//	client *redis.Client
//}
//
//func (r RedisClickCounterReader) GetTopToday(ctx context.Context, num int64) ([]int, error) {
//	key := keys.HotLinksKey + ":" + time.Now().Format("20060102")
//
//	if num <= 0 {
//		num = 100
//	}
//
//	res, err := r.client.ZRevRangeWithScores(ctx, key, 0, num).Result()
//	if err != nil {
//		return nil, fmt.Errorf("")
//	}
//
//}
