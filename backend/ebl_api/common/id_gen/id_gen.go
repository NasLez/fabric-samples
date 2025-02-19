package id_gen

import (
	"github.com/sony/sonyflake"
	"time"
)

var idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{
	StartTime: time.Now(),
})

func NextID() (int64, error) {
	id, err := idGenerator.NextID()
	return int64(id), err
}
