/*
 ____               _     _
|  _ \ ___ _ __ ___(_)___| |_ ___ _ __   ___ ___
| |_) / _ \ '__/ __| / __| __/ _ \ '_ \ / __/ _ \
|  __/  __/ |  \__ \ \__ \ ||  __/ | | | (_|  __/
|_|   \___|_|  |___/_|___/\__\___|_| |_|\___\___|

Copyright (C) Pocket52@2018 - All Rights Reserved

Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
*/
package redis_persistence

import (
	"gopkg.in/redis.v5"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var once sync.Once
var redisInstance *Redis

type Redis struct {
	serverSession *redis.Client
}

func GetRedis() *Redis {
	once.Do(func() {
		redisInstance = initialize()
	})

	return redisInstance
}

// init
func initialize() *Redis {
	// get redis_persistence host
	redisHost := getRedisHost()

	// wait for redis_persistence host to alive in 100 secs
	waitForRedis(redisHost, 10, 10)

	serverSession := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // no password set
		DB:       5,  // use default DB
	})

	pong, err := serverSession.Ping().Result()
	if err != nil {
		log.Panic(pong, err)
	}

	return &Redis{serverSession: serverSession}
}

// getRedisHost function gives the redis_persistence host to connect. Its generally configured from environment
// variable REDIS_HOST. In local test environment you don't need to define REDIS_HOST as it defaults
// to 'localhost'. In case of production server you can specify redis_persistence host or docker compose can pass
// host name to connect to.
func getRedisHost() string {
	// get the environment variable
	host := strings.Trim(os.Getenv("REDIS_HOST"), " ")

	// if the variable does not exist, return defaults
	if len(host) == 0 {
		return "localhost:6379"
	}

	// return the host name and default ports
	return host + ":6379"
}

// waitForRedis is a helper function that allows one to wait for redis_persistence to be alive.
// This function is helpful in case of server is started by docker and redis_persistence might not be alive
// then server tries to connect it.
func waitForRedis(host string, retries int, retryDelay int) {
	// create a new client
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// defer closing the client
	defer client.Close()

	// we ant to retry 10 times with 10 seconds delay
	var err error = nil

	// try connecting to redis_persistence
	for tries := 1; tries <= retries; tries++ {
		// ping redis_persistence
		_, err = client.Ping().Result()
		if err == nil {
			break
		} else {
			log.Printf("unable to connect redis_persistence at %s : try[%d], waiting for %d\n", host, tries, retryDelay)

			// connection failed, wait
			time.Sleep(time.Duration(retryDelay) * time.Second)
		}
	}

	// after all retries check error status, if error then connection failed
	if err != nil {
		log.Panicf("unable to connect redis_persistence at %s, ", err)
	}
}
