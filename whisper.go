package whisper

import "time"

func Wait(delay time.Duration) {
	time.Sleep(delay)
}

func Retry(fn func() error, times int, delay *time.Duration) error {
	var err error
	for i := 0; i < times; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		if delay != nil {
			Wait(*delay)
		}
	}
	return err
}
