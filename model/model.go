package model

import "time"

// Department represents an organizational department.
type Department struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:200;not null"`
	Code        string    `json:"code" gorm:"uniqueIndex;size:100;not null"`
	ParentID    uint      `json:"parent_id" gorm:"index;default:0;comment:0 means root"`
	Path        string    `json:"path" gorm:"size:500;comment:full path like /1/2/3/"`
	Level       int       `json:"level" gorm:"default:1;comment:depth level"`
	Sort        int       `json:"sort" gorm:"default:0;comment:sort order"`
	Leader      string    `json:"leader" gorm:"size:100;comment:department leader name"`
	Description string    `json:"description" gorm:"size:500"`
	Status      int       `json:"status" gorm:"default:1;comment:1=active,0=disabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Position represents a job position/title.
type Position struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:200;not null"`
	Code        string    `json:"code" gorm:"uniqueIndex;size:100;not null"`
	Level       int       `json:"level" gorm:"default:1;comment:position level/rank"`
	Description string    `json:"description" gorm:"size:500"`
	Status      int       `json:"status" gorm:"default:1;comment:1=active,0=disabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrgMember represents a user's membership in the organization.
type OrgMember struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       uint      `json:"user_id" gorm:"index;not null"`
	DepartmentID uint      `json:"department_id" gorm:"index;not null"`
	PositionID   uint      `json:"position_id" gorm:"index"`
	IsLeader     bool      `json:"is_leader" gorm:"default:false;comment:is department leader"`
	JoinedAt     time.Time `json:"joined_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
