package lock

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// Lock 锁对象
type Lock struct {
	resource string
	token    string
	conn     *redis.Client
	timeout  time.Duration
}

func (lock *Lock) key() string {
	return fmt.Sprintf("redislock:%s", lock.resource)
}

// lock 上锁
func (lock *Lock) lock() (bool, error) {
	seted, err := lock.conn.SetNX(lock.key(), lock.token, lock.timeout).Result()
	if err != nil {
		return false, err
	}
	return seted, nil
}

// Unlock 放弃锁
func (lock *Lock) Unlock() error {
	_, err := lock.conn.Del(lock.key()).Result()
	return err
}

// TryLockWithTimeout 尝试加锁
func TryLockWithTimeout(conn *redis.Client, resource string, token string, timeout time.Duration) (*Lock, bool, error) {
	lock := &Lock{conn: conn, resource: resource, token: token, timeout: timeout}
	ok, err := lock.lock()
	if err != nil {
		lock = nil
	}
	return lock, ok, err
}

// TryLock 尝试加锁
func TryLock(conn *redis.Client, resource string, token string) (*Lock, bool, error) {
	// 默认超时时间为 200 毫秒
	return TryLockWithTimeout(conn, resource, token, 200*time.Millisecond)
}

// WaitingGetLock 等待加锁
func WaitingGetLock(conn *redis.Client, resource string, token string) (*Lock, error) {
	for i := 0; i <= 300; i++ {
		lock, ok, err := TryLockWithTimeout(conn, resource, token, 200*time.Millisecond)
		if err != nil {
			return nil, err
		}
		if !ok {
			time.Sleep(1 * time.Millisecond)
			continue
		}
		return lock, nil
	}
	return nil, errors.New("failed to get lock")
}
