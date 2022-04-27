package domain

type FileInfo struct {
	Id       string
	Name     string
	MimeType string
	ParentId string
	Path     string
	Data     []byte
}
