package storage

import (
	"errors"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/9toon/kinu/config"
	"github.com/9toon/kinu/logger"
)

type Storage interface {
	Open() error

	Fetch(key string) (*Object, error)

	PutFromBlob(key string, image []byte, contentType string, metadata map[string]string) error
	Put(key string, imageFile io.ReadSeeker, contentType string, metadata map[string]string) error

	List(key string) ([]StorageItem, error)

	Move(from string, to string) error
}

type StorageItem interface {
	Key() string
	Filename() string
	Extension() string
	ImageSize() string
}

type Object struct {
	Body     []byte
	Metadata map[string]string
}

var (
	ErrImageNotFound = errors.New("not found requested image")
)

type ErrInvalidStorageOption struct {
	error
	Message string
}

func (e *ErrInvalidStorageOption) Error() string { return e.Message }

var (
	AvailableStorageTypes = []string{"S3", "File"}
	ErrUnknownStorage     = errors.New("specify unknown storage.")
	selectedStorageType   string
)

func init() {
	if config.BackwardCompatibleMode {
		selectedStorageType = "BackwardCompatibleS3"
	} else {
		selectedStorageType = os.Getenv("KINU_STORAGE_TYPE")
		if len(selectedStorageType) == 0 {
			panic("must specify KINU_STORAGE_TYPE system environment.")
		}

		var isAvailableStorageType bool
		for _, storageType := range AvailableStorageTypes {
			if selectedStorageType == storageType {
				isAvailableStorageType = true
			}
		}

		if !isAvailableStorageType {
			panic("unknown KINU_STORAGE_TYPE " + selectedStorageType + ".")
		}
	}

	logger.WithFields(logrus.Fields{
		"storage_type": selectedStorageType,
	}).Info("setup storage")
}

func Open() (Storage, error) {
	switch selectedStorageType {
	case "S3":
		return openS3Storage()
	case "File":
		return openFileStorage()
	case "BackwardCompatibleS3":
		return openBackwardCompatibleS3Storage()
	default:
		return nil, ErrUnknownStorage
	}
}
