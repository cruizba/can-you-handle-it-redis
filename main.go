package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	// Define command-line flags
	clusterNodes := flag.String("cluster", "", "Comma-separated list of Redis cluster node addresses")
	numberOfGoroutines := flag.Int("goroutines", 10, "Number of goroutines to run")
	password := flag.String("password", "", "Redis password")
	sleep := flag.String("sleep", "1s", "Sleep duration between writes")
	sleepDuration, err := time.ParseDuration(*sleep)
	if err != nil {
		fmt.Println("Error parsing sleep duration:", err)
		os.Exit(1)
	}

	// Parse the flags
	flag.Parse()

	if *clusterNodes == "" {
		fmt.Println("Please provide a comma-separated list of Redis cluster node addresses")
		os.Exit(1)
	}

	if *password == "" {
		fmt.Println("Please provide a Redis password")
	}

	// Split the cluster nodes string into a slice
	nodes := strings.Split(*clusterNodes, ",")

	// Create a Redis cluster client
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Password: *password,
		Addrs:    nodes,
	})

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	for i := 0; i < *numberOfGoroutines; i++ {
		wg.Add(1)
		go func(ctx context.Context, index int, rdb *redis.ClusterClient, wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				time.Sleep(sleepDuration)
				select {
				case <-ctx.Done():
					return
				default:
					randomKeyAndValue := fmt.Sprintf("%d-%d", index, time.Now().UnixNano())
					key := fmt.Sprintf("key-%s", randomKeyAndValue)
					value := fmt.Sprintf("value-%s", randomKeyAndValue)
					err := rdb.Set(ctx, key, value, 0).Err()
					if err != nil {
						fmt.Println("Error writing to Redis:", err)
						continue
					}
					fmt.Printf("gr-%d: Wrote key %s with value %s\n", index, key, value)

					val, err := rdb.Get(ctx, "key").Result()
					if err != nil {
						fmt.Println("Error reading from Redis:", err)
						continue
					}
					fmt.Printf("gr-%d: Read key %s with value %s\n", index, key, val)

					err = rdb.Del(ctx, key).Err()
					if err != nil {
						fmt.Println("Error deleting from Redis:", err)
						continue
					}
					fmt.Printf("gr-%d: Deleted key %s\n", index, key)
				}
			}
		}(ctx, i, rdb, &wg)
	}

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-sigChan
	fmt.Println("Shutting down...")
	cancel()  // Cancel the context, signaling all goroutines to stop
	wg.Wait() // Wait for all goroutines to finish
}
