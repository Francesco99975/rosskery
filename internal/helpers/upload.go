package helpers

import (
	"fmt"
	"image"
	_ "image/jpeg" // For JPEG format
	_ "image/png"  // For PNG format
	"io"
	"mime/multipart"
	"os"
	"path"

	"github.com/chai2010/webp"
)

func ImageUpload(file *multipart.FileHeader, topic string, identifier string) (string, error) {
	if file.Size > 1024*1024 {
		return "", fmt.Errorf("Image size cannot be greater than 1MB")
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	_, err = src.Seek(0, 0)
	if err != nil {
		return "", err
	}

	if !isFormatAllowed(buffer) {
		return "", fmt.Errorf("Image format not allowed")
	}

	_, err = os.Stat(path.Join("static", topic))
	if os.IsNotExist(err) {
		// Directory doesn't exist, create it
		err := os.MkdirAll(path.Join("static", topic), 0755)
		if err != nil {
			return "", fmt.Errorf("Error creating directory: %v", err)
		}

	}

	dst, err := os.Create(path.Join("static", topic, identifier+".webp"))
	if err != nil {
		return "", err
	}

	if !magicMatches(buffer, []byte{'R', 'I', 'F', 'F'}) {
		if err := ToWebP(&src, dst); err != nil {
			return "", err
		}
	} else {
		if _, err := io.Copy(dst, src); err != nil {
			return "", err
		}
	}
	defer dst.Close()

	return fmt.Sprintf("/static/%s/%s.webp", topic, identifier), nil
}

func DeleteImage(topic string, identifier string) error {
	return os.Remove(path.Join("static", topic, identifier+".webp"))
}

func ToWebP(file *multipart.File, uploadFile *os.File) error {
	// Decode the input image
	img, _, err := image.Decode(*file)
	if err != nil {
		return fmt.Errorf("unable to decode the image file: %s", err.Error())
	}

	// Encode the image to WebP format
	err = webp.Encode(uploadFile, img, nil)
	if err != nil {
		return fmt.Errorf("unable to encode the image to WebP format: %s", err.Error())
	}

	return err
}

func isFormatAllowed(buffer []byte) bool {
	jpegMagic := []byte{0xFF, 0xD8}
	pngMagic := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}
	webpMagic := []byte{'R', 'I', 'F', 'F'}

	return magicMatches(buffer, jpegMagic) || magicMatches(buffer, pngMagic) || magicMatches(buffer, webpMagic)
}

func magicMatches(buffer, magic []byte) bool {
	return len(buffer) >= len(magic) && bytesEqual(buffer[:len(magic)], magic)
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	} else {
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}

		return true
	}
}
