package dbmodel

import (
	"github.com/fahmifan/ulids"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type Base struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
	Metadata
}

type Metadata struct {
	CreatedAt null.Time
	UpdatedAt null.Time
	DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_at"`
}

type User struct {
	Base
	Name     string
	Email    string
	Password string
	Role     string
	Active   bool
}

func (user User) IsActive() bool {
	return user.Active
}

type FileExt string
type FileType string

const (
	FileTypeAssignmentCaseInput  FileType = "assignment_case_input"
	FileTypeAssignmentCaseOutput FileType = "assignment_case_output"
	FileTypeSubmission           FileType = "submission"
)

type File struct {
	Base
	Name string
	Type FileType
	Ext  FileExt
	URL  string
}

type OutboxItem struct {
	ID            ulids.ULID
	IdempotentKey string
	Status        string
	JobType       string
	Payload       string
	Version       int32
}
