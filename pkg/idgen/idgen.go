package idgen

import (
	"fmt"

	"github.com/sony/sonyflake"
)

var instance *sonyflake.Sonyflake

func init() {
	instance = sonyflake.NewSonyflake(sonyflake.Settings{})
}

func NextID() (uint64, error) {
	id, err := instance.NextID()
	if err != nil {
		return 0, fmt.Errorf("idgen: failed to generate id: %w", err)
	}
	return id, nil
}
