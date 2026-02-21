// helpers/file_uploader.go
package helpers

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type UploadResult struct {
	RelativePath string // Path untuk disimpan ke DB, contoh: "/uploads/transactions/abc.png"
	AbsolutePath string // Path fisik di server, untuk keperluan cleanup jika perlu
}

// SaveBase64Image menerima data URI dan menyimpannya ke server
// Returns: UploadResult dengan relative & absolute path, atau error
func SaveBase64Image(base64Data, uploadDir, baseURL string) (*UploadResult, error) {
	if base64Data == "" {
		return nil, fmt.Errorf("payment proof is required")
	}

	// 1. Validasi & extract MIME type + content
	if !strings.HasPrefix(base64Data, "data:image/") {
		return nil, fmt.Errorf("invalid image format: expected data URI")
	}

	parts := strings.SplitN(base64Data, ",", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid base64 data format")
	}

	header := parts[0]
	content := parts[1]

	// 2. Tentukan ekstensi file berdasarkan MIME type
	var ext string
	switch {
	case strings.Contains(header, "image/png"):
		ext = ".png"
	case strings.Contains(header, "image/jpeg") || strings.Contains(header, "image/jpg"):
		ext = ".jpg"
	default:
		return nil, fmt.Errorf("unsupported image type: only png, jpg, jpeg allowed")
	}

	// 3. Decode base64
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// 4. Validasi ukuran file (max 5MB)
	const maxSize = 5 << 20 // 5MB
	if len(decoded) > maxSize {
		return nil, fmt.Errorf("file too large: max 5MB allowed")
	}

	// 5. Generate unique filename
	filename := fmt.Sprintf("trx_%s%s", uuid.New().String(), ext)

	// 6. Pastikan directory upload ada
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// 7. Simpan file
	absolutePath := filepath.Join(uploadDir, filename)
	if err := os.WriteFile(absolutePath, decoded, 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	// 8. Return relative path untuk DB + absolute path untuk cleanup
	relativePath := filepath.ToSlash(filepath.Join(baseURL, filename))
	return &UploadResult{
		RelativePath: relativePath,
		AbsolutePath: absolutePath,
	}, nil
}

// CleanupUploadedFile menghapus file yang sudah terlanjur diupload (untuk rollback)
func CleanupUploadedFile(absolutePath string) {
	if absolutePath != "" {
		_ = os.Remove(absolutePath) // Ignore error, yang penting usaha cleanup
	}
}
