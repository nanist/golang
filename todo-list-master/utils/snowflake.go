package utils

import (
	"github.com/bwmarrin/snowflake"
	"time"
)

var node *snowflake.Node

func InitSnowflake(startTime string) (err error) {
	var begin time.Time
	begin, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	snowflake.Epoch = begin.UnixNano() / 1000000
	node, err = snowflake.NewNode(1)
	if err != nil {
		return
	}
	return
}
func GenID() int64 {
	return node.Generate().Int64()
}
