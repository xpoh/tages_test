package filestorage

import (
	"os"
	"time"
)

type FileInStorageInfo struct {
	user       string `json:"User"`
	Name       string `json:"Name"`
	Path       string `json:"Path"`
	CreateDate string `json:"CreateDate"`
	ModifyDate string `json:"ModifyDate"`
}

type FileStorager interface {
	GetFileList(user string) ([]FileInStorageInfo, error)
	GetFile(user string, filename string) ([]byte, error)
	PutFile(user string, filename string, file []byte) error
}

type ImMemoryLocalStorage struct {
	files []FileInStorageInfo
}
type FileNotFoundError struct{}

func (f FileNotFoundError) Error() string {
	return "File not found"
}

type WriteFileError struct{}

func (f WriteFileError) Error() string {
	return "WriteFileError"
}

func (i *ImMemoryLocalStorage) PutFile(user string, filename string, file []byte) error {
	if i.files == nil {
		i.files = make([]FileInStorageInfo, 0)
	}
	if err := os.WriteFile(user+"-"+filename, file, os.ModePerm); err != nil {
		return err
	}
	i.files = append(i.files, FileInStorageInfo{
		user:       user,
		Name:       filename,
		Path:       user + "-" + filename,
		CreateDate: time.Now().GoString(),
		ModifyDate: time.Now().GoString(),
	})
	return nil
}

func (i *ImMemoryLocalStorage) GetFile(user string, filename string) ([]byte, error) {
	for _, file := range i.files {
		if file.user == user && file.Name == filename {
			if f, err := os.ReadFile(file.Path); err == nil {
				return f, nil
			}
		}
	}
	return nil, FileNotFoundError{}
}

func (i *ImMemoryLocalStorage) GetFileList(user string) ([]FileInStorageInfo, error) {
	result := make([]FileInStorageInfo, 0, len(i.files))
	for _, v := range i.files {
		if v.user == user {
			result = append(result, v)
		}
	}

	return result, nil
}
