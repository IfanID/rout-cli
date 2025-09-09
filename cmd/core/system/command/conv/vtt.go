package conv

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func ParseVTTDuration(ts string) (time.Duration, error) {
	var h, m, s, ms int
	parts := strings.Split(ts, ":")
	if len(parts) == 3 {
		h, _ = strconv.Atoi(parts[0])
		m, _ = strconv.Atoi(parts[1])
		sMs := strings.Split(parts[2], ".")
		s, _ = strconv.Atoi(sMs[0])
		ms, _ = strconv.Atoi(sMs[1])
	} else if len(parts) == 2 {
		m, _ = strconv.Atoi(parts[0])
		sMs := strings.Split(parts[1], ".")
		s, _ = strconv.Atoi(sMs[0])
		ms, _ = strconv.Atoi(sMs[1])
	} else {
		return 0, fmt.Errorf("format timestamp tidak valid: %s", ts)
	}
	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second + time.Duration(ms)*time.Millisecond, nil
}

func FormatVTTDuration(d time.Duration) string {
	d = d.Round(time.Millisecond)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	ms := d / time.Millisecond
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}

func ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func InsertStringSlice(slice []string, index int, values ...string) []string {
	if index < 0 || index > len(slice) {
		return slice
	}
	result := make([]string, len(slice)+len(values))
	copy(result[:index], slice[:index])
	copy(result[index:index+len(values)], values)
	copy(result[index+len(values):], slice[index:])
	return result
}
