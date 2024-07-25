package models

import (
	"strings"
	"errors"

	log "github.com/gophish/gophish/logger"
)

type File struct {
	Id	int64		`json:"id"`
	Name	string		`json:"name"`
	Parent	*Directory	`json:"parent"`
	Size	int64		`json:"size"`
}

type Directory struct {
	Id	int64		`json:"id"`
	Name	string		`json:"name"`
	Path	string		`json:"path"`
	Parent	*Directory	`json:"parent"`
	Subdirs	[]Directory	`json:"subdirs"`
	Files	[]File		`json:"files"`
}

// ErrInvalidPath is thrown to block potential path traversal attacks
var ErrInvalidPath = errors.New("Path blocked")

// ErrInvalidParent is thrown when a file or directory is being 
// created with an invalid parent directory
var ErrInvalidParent = errors.New("Invalid parent")

// ErrFileNotFound is thrown when a non-existent file is being accessed
var ErrFileNotFound = errors.New("File not found")

// ErrDirectoryNotFound is thrown when a non-existent directory is being accessed
var ErrDirectoryNotFound = errors.New("Directory not found")

// Validate performs validation on a file
func (f *File) Validate() error {
	return nil
}

// Validate performs validation on a directory
func (d *Directory) Validate() error {
	if strings.Contains(d.Path, "..") {
		return ErrInvalidPath
	}
	return nil
}

func GetSubdirs(id int64) ([]Directory, error) {
	subdirs := []Directory{}
	err := db.Where("parent=?", id).Find(&subdirs).Error
	if err != nil {
		log.Error(err)
	}
	return subdirs, err
}

func GetFiles(id int64) ([]File, error) {
	files := []File{}
	err := db.Where("parent=?", id).Find(&files).Error
	if err != nil {
		log.Error(err)
	}
	return files, err
}

func GetDirectory(id int64) (Directory, error) {
	dir := Directory{}
	err := db.Where("id=?", id).Find(&dir).Error
	if err != nil {
		log.Error(err)
		return dir, err
	}
	dir.Subdirs, err = GetSubdirs(dir.Id)
	if err != nil {
		log.Error(err)
		return dir, err
	}
	dir.Files, err = GetFiles(dir.Id)
	if err != nil {
		log.Error(err)
	}
	return dir, err
}

func UploadFile(f *File) error {
	if err := f.Validate(); err != nil {
		return err
	}
	tx := db.Begin()
	err := tx.Save(f).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
	}
	return err
}

func CreateDirectory(d *Directory) error {
	if err := d.Validate(); err != nil {
		return err
	}
	tx := db.Begin()
	err := tx.Save(d).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
	}
	return err
}
