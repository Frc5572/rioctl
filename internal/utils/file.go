package utils

type File struct {
	Size    int64
	Current int64
	Name    string
}

func (m File) UpdateCurrent(current int64) File {
	m.Current = current
	return m
}

func NewFile(size int64, current int64, name string) File {
	return File{
		Name:    name,
		Current: current,
		Size:    size,
	}
}
