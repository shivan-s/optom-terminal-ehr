package db

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// TODO: create custom "titles"

type Patient struct {
	gorm.Model
	FirstName  string
	LastName   string
	MiddleName string
	Title      string
	BirthDate  datatypes.Date
	Address    string
	EyeExams   []EyeExam
}

type PatientRepository struct {
	db *gorm.DB
	pt []*Patient
}

// AllPatients gives list of all patients. Will probably not use this.
func (p *PatientRepository) AllPatients() (patients []*Patient, err error) {
	err = p.db.Order("lastname").Find(&p.pt).Error
	return patients, err
}

// FindPatient gives list of patients based on a search term, `s`, of which
// will search first, middle and lastname.
func (p *PatientRepository) FindPatients(s string) (patients []*Patient, err error) {
	// TODO does SOUNDEX exist? try this.
	err = p.db.Where("firstname LIKE ? OR lastname LIKE ? OR middlename LIKE ?", s, s, s).Find(&patients).Error
	return patients, err
}

// PatientById gives a single patient by id.
func (p *PatientRepository) PatientById(id uint) (patient *Patient, err error) {
	err = p.db.First(&p.pt, id).Error
	return patient, err
}

// CreatePatient creates a patient.
func (p *PatientRepository) CreatePatient(patient *Patient) (err error) {
	err = p.db.Create(&p.pt).Error
	return err
}

// UpdatePatient updates a patient.
func (p *PatientRepository) UpdatePatient(patient *Patient) (err error) {
	// TODO: need to write test
	err = p.db.Model(&p.pt).Updates(patient).Error
	return err
}

// DeletePatient deletes a patient.
func (p *PatientRepository) DeletePatient(id uint) (err error) {
	// TODO: need to write test
	err = p.db.Delete(&p.pt, id).Error
	return err
}
