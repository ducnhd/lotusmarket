package fetchers

import (
	"fmt"
	"strconv"
	"strings"
)

func getString(m map[string]any, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func getFloat(m map[string]any, key string) float64 {
	if v, ok := m[key]; ok && v != nil {
		switch n := v.(type) {
		case float64:
			return n
		case string:
			f, _ := strconv.ParseFloat(n, 64)
			return f
		}
	}
	return 0
}

func getInt(m map[string]any, key string) int64 { return int64(getFloat(m, key)) }

func parseFloatStr(s string) float64 { f, _ := strconv.ParseFloat(s, 64); return f }

func parsePipeVol(s string) int64 {
	parts := strings.Split(s, "|")
	if len(parts) >= 2 {
		v, _ := strconv.ParseInt(parts[1], 10, 64)
		return v
	}
	return 0
}
