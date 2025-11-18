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
