package meta

import (
	"io"
)

type Application struct {
	ID   string
	Data []byte
}

func readApplication(reader io.Reader, size int) (*Application, error) {
	application := &Application{}

	// 4 bytes per ID
	id := make([]byte, 4)
	_, err := reader.Read(id)
	if err != nil {
		return application, err
	}
	application.ID = string(id)

	// all another data for application
	data := make([]byte, size-4)
	_, err = reader.Read(data)
	if err != nil {
		return application, err
	}
	application.Data = data

	return application, nil
}
