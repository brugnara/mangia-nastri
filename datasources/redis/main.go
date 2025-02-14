package redis

import (
	"context"
	"encoding/json"
	"mangia_nastri/datasources"
	"mangia_nastri/logger"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDataSource struct {
	log         logger.Logger
	redisClient *redis.Client
}

var ctx = context.Background()

func New(log *logger.Logger, connectionUrl string) *RedisDataSource {
	options, err := redis.ParseURL(connectionUrl)
	if err != nil {
		log.Fatalf("Error parsing connection URL: %v", err)
	}

	return &RedisDataSource{
		log:         log.CloneWithPrefix("redis"),
		redisClient: redis.NewClient(options),
	}
}

func (ds *RedisDataSource) Ready() <-chan bool {
	// do a PING and check the response is PONG
	ch := make(chan bool)
	wg := sync.WaitGroup{}

	// try 5 times then fail
	go func() {
		connected := false
		for i := 0; i < 5 && !connected; i++ {
			wg.Add(1)

			go func() {
				_, err := ds.redisClient.Ping(ctx).Result()

				wg.Done()

				if err != nil {
					ds.log.Printf("Error pinging Redis: %v", err)
					return
				}

				connected = true
			}()

			wg.Wait()

			time.Sleep(1 * time.Second)
		}

		if !connected {
			ds.log.Fatalf("Redis is not ready after 5 attempts")
		} else {
			ds.log.Info("Redis is ready")
			close(ch)
		}
	}()

	return ch
}

func (ds *RedisDataSource) Set(key datasources.Hash, value datasources.Payload) error {
	jsonString, err := json.Marshal(value)

	if err != nil {
		ds.log.Printf("SET: Error marshalling value for key %v: %v", key, err)
		return err
	}

	result := ds.redisClient.Set(ctx, key.String(), jsonString, 0)

	if result.Err() != nil {
		ds.log.Printf("SET: Error setting value for key %v: %v", key, result.Err())
		return result.Err()
	}

	ds.log.Printf("SET: Setting value %v for key %v", value, key)

	return nil
}

func (ds *RedisDataSource) Get(key datasources.Hash) (datasources.Payload, error) {
	result := ds.redisClient.Get(ctx, key.String())

	if result.Err() != nil {
		ds.log.Printf("GET: Error getting value for key %v: %v", key, result.Err())
		return datasources.Payload{}, result.Err()
	}

	var payload datasources.Payload
	err := json.Unmarshal([]byte(result.Val()), &payload)

	if err != nil {
		ds.log.Printf("GET: Error unmarshalling value for key %v: %v", key, err)
		return datasources.Payload{}, err
	}

	ds.log.Printf("GET: Getting value for key %v", key)

	return payload, nil
}
