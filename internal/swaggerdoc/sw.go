package swaggerdoc

import "os"

type Sw struct {
	File string
}

func (s *Sw) ReadDoc() string {
	bytes, err := os.ReadFile(s.File)
	if err != nil {
		return ""
	}
	return string(bytes)
}
