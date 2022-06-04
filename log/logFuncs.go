package log

import (
	"fmt"
	"time"
)

func LogNow() {
	fmt.Println(time.Now().Format(" 2006-01-02 15:04:05"))
}
