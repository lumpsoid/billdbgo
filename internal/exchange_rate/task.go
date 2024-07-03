package exchangerate

import "time"

func scheduleTasks() {
	// Define a ticker that triggers every minute
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Run task 1 every minute
			// go task1()

			// Run task 2 every 5 minutes
			// if time.Now().Minute()%5 == 0 {
			// go task2()
			// }
		}
	}
}
