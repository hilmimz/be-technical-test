package redis

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type StockRepository struct {
	client *redis.Client
}

func NewStockRepository(client *redis.Client) *StockRepository {
	return &StockRepository{client: client}
}

func stockKey(sku string) string {
	return "stock:" + sku
}

func (r *StockRepository) SetStock(ctx context.Context, sku string, qty int) error {
	return r.client.Set(ctx, stockKey(sku), qty, 0).Err()
}

var decrementScript = redis.NewScript(`
	local stock = redis.call('GET', KEYS[1])
	if stock == false then
		return -1
	end
	stock = tonumber(stock)
	local qty = tonumber(ARGV[1])
	if stock >= qty then
		redis.call('DECRBY', KEYS[1], qty)
		return 1
	else
		return 0
	end
`)

func (r *StockRepository) DecrementStock(ctx context.Context, sku string, qty int) (int, error) {
	result, err := decrementScript.Run(ctx, r.client, []string{stockKey(sku)}, qty).Result()
	if err != nil {
		return 0, err
	}

	val, ok := result.(int64)
	if !ok {
		parsed, parseErr := strconv.Atoi(result.(string))
		if parseErr != nil {
			return 0, parseErr
		}
		return parsed, nil
	}
	return int(val), nil
}

func (r *StockRepository) IncrementStock(ctx context.Context, sku string, qty int) error {
	return r.client.IncrBy(ctx, stockKey(sku), int64(qty)).Err()
}
