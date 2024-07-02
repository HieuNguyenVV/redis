package main

import (
	"context"
	"fmt"
	"log"
	"redis/cache"
	"time"
)

func main() {
	duration := time.Millisecond * 1000
	ctx := context.Background()
	redisClient, err := cache.NewRedis("localhost:6379", "", "", duration, duration)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	exist, err := redisClient.Exist(ctx, "key_1")
	if err != nil {
		fmt.Println(err)
		return
	}
	if exist != 1 {
		if err := redisClient.Set(ctx, "key_1", "VALUE_KEY_1", duration); err != nil {
			fmt.Println(err)
			return
		}
	}

	value, err := redisClient.Get(ctx, "key_1")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(value)

	entries := []cache.Entry{
		{Key: "key_2", Value: "VALUE_KEY_2"},
		{Key: "key_3", Value: "VALUE_KEY_3"},
		{Key: "key_4", Value: "VALUE_KEY_4"},
	}

	if err := redisClient.MSet(ctx, entries, duration); err != nil {
		fmt.Println(err)
		return
	}

	results, err := redisClient.MGet(ctx, []string{"key_2", "key_3", "key_4"})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, value := range results {
		fmt.Println(value)
	}

	if err := redisClient.Del(ctx, "key_1"); err != nil {
		fmt.Println(err)
		return
	}
	value, err = redisClient.Get(ctx, "key_1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(value)

	if err := redisClient.MDel(ctx, []string{"key_2", "key_3", "key_4"}); err != nil {
		fmt.Println(err)
		return
	}
	results, err = redisClient.MGet(ctx, []string{"key_2", "key_3", "key_4"})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, value := range results {
		fmt.Println(value)
	}
}
