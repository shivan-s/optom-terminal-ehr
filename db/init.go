package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type Optometrist struct {
	gorm.Model
	Name           string
	BoardRegNumber string
}

type Dispenser struct {
	gorm.Model
	Name           string
	BoardRegNumber string
}

type EyeExam struct {
	gorm.Model
	History       string
	Health        string
	Management    string
	ExamAt        time.Time
	Prescriptions []Prescription
	OptometristID uint
	PatientID     uint
}

type Prescription struct {
	gorm.Model
	RightSphere          float32
	RightCylinder        float32
	RightAxis            float32
	RightInterAdd        float32
	RightAdd             float32
	RightHorizontalPrism float32
	RightVerticalPrism   float32
	LeftSphere           float32
	LeftCylinder         float32
	LeftAxis             float32
	LeftInterAdd         float32
	LeftAdd              float32
	LeftHorizontalPrism  float32
	LeftVerticalPrism    float32
	OptometristID        uint
	PatientID            uint
	EyeExamID            uint
}

type Job struct {
	gorm.Model
	Frame          string
	PrescriptionID uint
	PatientID      uint
	DispenserID    uint
}

// SetUpDb runs migrations on the database.
func SetUpDb() (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open("../optom.db"), &gorm.Config{})
	err = db.AutoMigrate(
		&Optometrist{},
		&Dispenser{},
		&Patient{},
		&EyeExam{},
		&Prescription{},
		&Job{},
	)
	return db, err
}
