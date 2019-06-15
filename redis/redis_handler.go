package redis_persistence

import (
	"encoding/json"
	"fmt"
	"log"
)

// Add table stat to redis_persistence
func (redis *Redis) InsertTailData(data interface{}) error {
	dataJson, _ := json.Marshal(data)
	dataJsonString := string(dataJson)

	pipe := redis.serverSession.Pipeline()
	pipe.RPush("client", dataJsonString)
	_, err := pipe.Exec()

	if err != nil {
		log.Println("[ERROR] LPUSH ", err)
	} else {
		log.Println("[SUCCESS] LPUSH ", data)
	}
	pipe.Close()

	return err
}

func (redis *Redis) ReadHeadData() string {
	pipe := redis.serverSession
	data := pipe.LRange("client", 0, 0).Val()

	fmt.Println(data)
	if len(data) > 0 {
		return data[0]
	}

	return ""
}

func (redis *Redis) PopHeadData() string {
	pipe := redis.serverSession
	data := pipe.BLPop(10, "client").Val()

	if len(data) > 0 {
		fmt.Println(data)
		return data[0]
	}

	return ""
}
