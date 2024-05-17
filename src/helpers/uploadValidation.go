package helpers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SizeUploadValidation(fileSize int64, maxFileSize int64) error {
	if fileSize > maxFileSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "Ukuran file melebihi 2MB")
	}
	return nil
}

func TypeUploadValidation(buffer []byte, validFileTypes []string) error {
	fileType := http.DetectContentType(buffer)
	if !isValidFileType(validFileTypes, fileType) {
		return fiber.NewError(fiber.StatusBadRequest, "Tipe file tidak valid. Hanya png, jpg, jpeg, dan pdf yang diperbolehkan.")
	}
	return nil
}

func isValidFileType(validFileTypes []string, fileType string) bool {
	for _, validType := range validFileTypes {
		if validType == fileType {
			return true
		}
	}
	return false
}
