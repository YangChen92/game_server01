package utils

import (
	"fmt"
	"sync"
	"time"
)

// Snowflake 结构体
type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	datacenterID  int64
	sequence      int64
}

// 定义常量
const (
	workerIDBits     = int64(5)
	datacenterIDBits = int64(5)
	sequenceBits     = int64(12)

	maxWorkerID     = int64(-1) ^ (int64(-1) << workerIDBits)
	maxDatacenterID = int64(-1) ^ (int64(-1) << datacenterIDBits)
	sequenceMask    = int64(-1) ^ (int64(-1) << sequenceBits)

	workerIDShift      = sequenceBits
	datacenterIDShift  = sequenceBits + workerIDBits
	timestampLeftShift = sequenceBits + workerIDBits + datacenterIDBits
	epoch              = int64(1609459200000) // 自定义起始时间(2021-01-01)
)

// 创建Snowflake实例
func NewSnowflake(workerID, datacenterID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("worker ID must be between 0 and %d", maxWorkerID)
	}
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, fmt.Errorf("datacenter ID must be between 0 and %d", maxDatacenterID)
	}
	return &Snowflake{
		lastTimestamp: 0,
		workerID:      workerID,
		datacenterID:  datacenterID,
		sequence:      0,
	}, nil
}

// 生成ID
func (s *Snowflake) GenerateID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixNano() / 1000000

	if timestamp < s.lastTimestamp {
		panic("clock moved backwards")
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	return ((timestamp - epoch) << timestampLeftShift) |
		(s.datacenterID << datacenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence
}
