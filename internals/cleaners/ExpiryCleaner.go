package cleaners

import (
	"godis/internals/handlers"
	"time"
)

func StartExpiryCleaner() {
	go func() {
		for {
			now := time.Now().Unix()
			handlers.TTLsMS.Lock()
			for key, expiredAt := range handlers.TTLs {
				if expiredAt <= now {
					delete(handlers.TTLs, key)
					handlers.SETsMS.Lock()
					delete(handlers.SETs, key)
					handlers.SETsMS.Unlock()

				}
			}
			handlers.TTLsMS.Unlock()

			time.Sleep(1 * time.Second)
		}
	}()
}

func Expired(key string) bool {
	handlers.TTLsMS.RLock()
	ts, ok := handlers.TTLs[key]
	handlers.TTLsMS.Unlock()

	if !ok {
		return false
	}

	if ts <= time.Now().Unix() {
		handlers.TTLsMS.Lock()
		delete(handlers.TTLs, key)
		handlers.TTLsMS.Unlock()

		handlers.SETsMS.Lock()
		delete(handlers.SETs, key)
		handlers.SETsMS.Unlock()
		return true
	}

	return false
}
