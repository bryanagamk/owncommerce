package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type LocalStorage struct {
	rootPath string
}

func NewLocalStorage(rootPath string) (*LocalStorage, error) {
	abs, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return nil, err
	}
	return &LocalStorage{rootPath: abs}, nil
}

func (s *LocalStorage) TenantDir(tenantID uuid.UUID) (string, error) {
	dir := filepath.Join(s.rootPath, "tenants", tenantID.String())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func (s *LocalStorage) Save(tenantID uuid.UUID, category, filename string, reader io.Reader) (string, error) {
	base, err := s.TenantDir(tenantID)
	if err != nil {
		return "", err
	}

	category = sanitizePath(category)
	filename = sanitizeFilename(filename)
	dir := filepath.Join(base, category)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	fullPath := filepath.Join(dir, filename)
	file, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", err
	}

	rel := filepath.ToSlash(filepath.Join("tenants", tenantID.String(), category, filename))
	return rel, nil
}

func (s *LocalStorage) Root() string {
	return s.rootPath
}

func sanitizePath(p string) string {
	p = strings.Trim(p, "/\\")
	p = strings.ReplaceAll(p, "..", "")
	if p == "" {
		return "misc"
	}
	return p
}

func sanitizeFilename(name string) string {
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, "..", "")
	if name == "" {
		return uuid.New().String()
	}
	return name
}

func (s *LocalStorage) FullPath(relativePath string) (string, error) {
	rel := filepath.Clean(relativePath)
	if strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("invalid path")
	}
	return filepath.Join(s.rootPath, rel), nil
}
