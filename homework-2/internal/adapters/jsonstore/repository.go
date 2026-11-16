package jsonstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"pos/internal/application"
	"pos/internal/domain"
)

type Repository struct {
	directory string
}

func NewRepository(directory string) (*Repository, error) {
	if strings.TrimSpace(directory) == "" {
		return nil, errors.New("storage directory must not be empty")
	}
	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}
	return &Repository{directory: directory}, nil
}

func (r *Repository) coursePath(id domain.CourseID) string {
	return filepath.Join(r.directory, string(id)+".json")
}

func (r *Repository) Load(id domain.CourseID) (*domain.Course, error) {
	validated, err := domain.NewCourse(id)
	if err != nil {
		return nil, err
	}
	id = validated.ID()
	content, err := os.ReadFile(r.coursePath(id))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: %q", application.ErrCourseNotFound, id)
		}
		return nil, fmt.Errorf("read course file: %w", err)
	}
	var data courseData
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("decode course file: %w", err)
	}
	var extra json.RawMessage
	if err := decoder.Decode(&extra); err != io.EOF {
		if err != nil {
			return nil, fmt.Errorf("read JSON data: %w", err)
		}
		return nil, errors.New("course file contains multiple JSON values")
	}
	course, err := courseFromData(data)
	if err != nil {
		return nil, fmt.Errorf("restore course from file: %w", err)
	}
	if course.ID() != id {
		return nil, fmt.Errorf("course ID mismatch: requested %q, file contains %q", id, course.ID())
	}
	return course, nil
}

func (r *Repository) Save(course *domain.Course) error {
	data, err := courseToData(course)
	if err != nil {
		return fmt.Errorf("prepare course for saving: %w", err)
	}
	if _, err := domain.NewCourse(course.ID()); err != nil {
		return fmt.Errorf("validate course ID: %w", err)
	}
	content, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("encode course: %w", err)
	}
	content = append(content, '\n')
	file, err := os.CreateTemp(r.directory, ".course-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary course file: %w", err)
	}
	tempPath := file.Name()
	defer os.Remove(tempPath)
	defer file.Close()
	if _, err := file.Write(content); err != nil {
		return fmt.Errorf("write course file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close course file: %w", err)
	}
	if err := os.Rename(tempPath, r.coursePath(course.ID())); err != nil {
		return fmt.Errorf("replace course file: %w", err)
	}
	return nil
}

func (r *Repository) ListIDs() ([]domain.CourseID, error) {
	entries, err := os.ReadDir(r.directory)
	if err != nil {
		return nil, fmt.Errorf("read storage directory: %w", err)
	}
	ids := make([]domain.CourseID, 0)
	for _, entry := range entries {
		if !entry.Type().IsRegular() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		id := domain.CourseID(strings.TrimSuffix(entry.Name(), ".json"))
		if strings.TrimSpace(string(id)) != string(id) {
			return nil, fmt.Errorf("invalid course filename: %q", entry.Name())
		}
		if _, err := r.Load(id); err != nil {
			return nil, fmt.Errorf("list course file %q: %w", entry.Name(), err)
		}
		ids = append(ids, id)
	}
	slices.Sort(ids)
	return ids, nil
}

var _ application.CourseRepository = (*Repository)(nil)
