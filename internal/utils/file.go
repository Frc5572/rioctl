package utils

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
	Size string
	Name string
}
