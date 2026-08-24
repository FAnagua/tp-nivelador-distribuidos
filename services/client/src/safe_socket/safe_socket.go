package safe_socket

import (
	"io"
)

//TODO: Complete with a short-read/short-write tolerant implementation

func SendAll(socket io.Writer, bytes []byte) error {
	totalBytes := len(bytes)

	for totalBytes > 0 {
		bytesWrite, err := socket.Write(bytes)
		if err != nil {
			return err
		}
		totalBytes -= bytesWrite
		bytes = bytes[bytesWrite:]
	}
	return nil
}

func RecvAll(socket io.Reader, size int) ([]byte, error) {
	data := make([]byte, 0, size)
	buff := make([]byte, size)
	for len(data) < size {
		bytesRead, err := socket.Read(buff[:size-len(data)])
		if err != nil {
			return nil, err
		}
		data = append(data, buff[:bytesRead]...)
	}
	return data, nil
}
