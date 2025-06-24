package model

import "time"

type Applicant struct {
    ApplicationID int        `db:"application_id"`
    Name          string     `db:"name"`
    Email         string     `db:"email"`
    Phone         *string    `db:"phone"`
    Education     *string    `db:"education"`
    Experience    *string    `db:"experience"`
    TechStack     *string    `db:"tech_stack"`
    CurrentStage  string     `db:"current_stage"`
    CreatedAt     time.Time  `db:"created_at"`
    CreatedBy     string     `db:"created_by"`
    UpdatedAt     *time.Time `db:"updated_at"`
    UpdatedBy     *string    `db:"updated_by"`
}

func (Applicant) TableName() string {
    return "applicant"
}
