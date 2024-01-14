package message

type ParsedFile struct {
	data interface{}
}

func NewParsedFile(data interface{}) ParsedFile {
	return ParsedFile{data: data}
}

func (msg ParsedFile) Data() interface{} {
	return msg.data
}
