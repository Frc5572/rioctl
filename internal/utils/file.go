package utils

import "fmt"

type FileUpload struct {
	File File
	// Size    int64
	Current int64
	// Name    string
}

func (m FileUpload) UpdateCurrent(current int64) FileUpload {
	m.Current = current
	return m
}

func NewFile(size int64, current int64, name string) FileUpload {
	return FileUpload{
		File: File{
			Name: name,
			Size: size,
		},
		Current: current,
	}
}

type File struct {
	Size      int64
	Name      string
	Event     string
	MatchType string
}

func (f File) HumanizedSize() string {
	return Humanize(f.Size)
}

func Humanize(b int64) string {
	return fmt.Sprintf("%.2f MB", float64(b)/1024/1024)
}
