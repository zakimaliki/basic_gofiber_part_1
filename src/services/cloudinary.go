package services

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gofiber/fiber/v2"
)

func UploadCloudinary(c *fiber.Ctx, file *multipart.FileHeader) (*uploader.UploadResult, error) {
	// Menginisialisasi konfigurasi Cloudinary
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Cloudinary URL not found")
	}

	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Membuka file
	src, err := file.Open()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	// Menghapus ekstensi dari file.Filename
	ext := filepath.Ext(file.Filename)
	fileNameWithoutExt := file.Filename[:len(file.Filename)-len(ext)]

	// Konfigurasi uploader Cloudinary
	uploadParams := uploader.UploadParams{
		PublicID:  fmt.Sprintf("%d_%s", time.Now().Unix(), fileNameWithoutExt),
		Overwrite: true,
	}

	// Mengunggah file ke Cloudinary
	uploadResult, err := cld.Upload.Upload(c.Context(), src, uploadParams)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return uploadResult, nil
}
