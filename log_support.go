package groschen

import (
	"fmt"
	"time"
)

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < len(runes)/2; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func FormatIntWithThousandSeparator(v int, sep string) string {
	tmp := reverse(fmt.Sprintf("%d", v))
	result := ""
	for len(tmp) > 3 {
		result += tmp[0:3] + sep
		tmp = tmp[3:]
	}
	result += tmp
	return reverse(result)
}

func FormatSpeed(bytes int, duration time.Duration) string {
	if duration.Seconds() < 0.1 {
		return "--.- K/s"
	}
	KBPS := (float64(bytes) / duration.Seconds()) / 1024.0
	return fmt.Sprintf("%.1f K/s", KBPS)
}

func FormatBytes(bytes int) string {
	return fmt.Sprintf("%s bytes", FormatIntWithThousandSeparator(bytes, ","))
}

const (
	LogStart = iota
	LogEnd   = iota
	LogOther = iota
)

type LogFunc func(logType int, prefix string, format string, args ...interface{})

func SeqLog(logType int, prefix string, format string, args ...interface{}) {
	fmt.Printf("%s: ", prefix)
	fmt.Printf(format, args...)
	fmt.Println()
}
