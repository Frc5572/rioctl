package utils

import "fmt"

type FileUpload struct {
	Size    int64
	Current int64
	Name    string
}

func (m FileUpload) UpdateCurrent(current int64) FileUpload {
	m.Current = current
	return m
}

func NewFile(size int64, current int64, name string) FileUpload {
	return FileUpload{
		Name:    name,
		Current: current,
		Size:    size,
	}
}

type File struct {
	Size      int64
	Name      string
	Event     string
	MatchType string
}

func Humanize(b int64) string {
	return fmt.Sprintf("%.2f MB", float64(b)/1024/1024)
}
