package db

import (
	// "database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/datatypes"
	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"regexp"
	"time"
)

var _ = Describe("Patient", func() {
	var repository *PatientRepository
	var mock sqlmock.Sqlmock

	BeforeEach(func() {
		var err error

		_, mock, err = sqlmock.New()
		Expect(err).ShouldNot(HaveOccurred())

		gdb, err := SetUpDb()
		Expect(err).ShouldNot(HaveOccurred())

		repository = &PatientRepository{db: gdb}

	})

	AfterEach(func() {
		err := mock.ExpectationsWereMet()
		Expect(err).ShouldNot(HaveOccurred())

	})

	Context("all", func() {
		It("all", func() {
			patient_one := &Patient{
				FirstName:  "One",
				LastName:   "PatientOne",
				MiddleName: "First",
				Title:      "Mr",
				BirthDate:  datatypes.Date(time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)),
				Address:    "One Address",
			}
			patient_two := &Patient{
				FirstName:  "Two",
				LastName:   "PatientTwo",
				MiddleName: "Second",
				Title:      "Mrs",
				BirthDate:  datatypes.Date(time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC)),
				Address:    "Two Address",
			}
			patient_three := &Patient{
				FirstName:  "Three",
				LastName:   "PatientThree",
				MiddleName: "Three",
				Title:      "Miss",
				BirthDate:  datatypes.Date(time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC)),
				Address:    "Three Address",
			}
			rows := sqlmock.NewRows(
				[]string{
					"firstname",
					"lastname",
					"middlename",
					"title",
					"birthdate",
					"address",
				}).AddRow(
				patient_one.ID,
				patient_one.FirstName,
				patient_one.LastName,
				patient_one.MiddleName,
				patient_one.Title,
				patient_one.BirthDate,
				patient_one.Address,
			).AddRow(
				patient_two.ID,
				patient_two.FirstName,
				patient_two.LastName,
				patient_two.MiddleName,
				patient_two.Title,
				patient_two.BirthDate,
				patient_two.Address,
			).AddRow(
				patient_three.ID,
				patient_three.FirstName,
				patient_three.LastName,
				patient_three.MiddleName,
				patient_three.Title,
				patient_three.BirthDate,
				patient_three.Address,
			)
			const sqlSelectAll = `SELECT * FROM "patients"`
			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectAll)).
				WillReturnRows(rows)

			result, err := repository.AllPatients()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result).Should(Equal(map[string]*Patient{
				"patient_one":   patient_one,
				"patient_two":   patient_two,
				"patient_three": patient_three,
			}))

			list, err := repository.AllPatients()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(list).Should(BeEmpty())

		})
		It("empty", func() {
			const sqlSelectAll = `SELECT * FROM "patients"`
			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectAll)).
				WillReturnRows(sqlmock.NewRows(nil))

			list, err := repository.AllPatients()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(list).Should(BeEmpty())
		})
	})

	Context("find", func() {
		It("found", func() {
			const searchTerm = "On"
			const sqlFind = `
        SELECT * FROM "patients"
        WHERE firstname LIKE ($1) OR
        WHERE lastname LIKE ($1) OR
        WHERE middlename LIKE ($1)
        `
			patient := &Patient{
				FirstName:  "One",
				LastName:   "PatientOne",
				MiddleName: "First",
				Title:      "Mr",
				BirthDate:  datatypes.Date(time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)),
				Address:    "One Address",
			}
			rows := sqlmock.NewRows(
				[]string{
					"firstname",
					"lastname",
					"middlename",
					"title",
					"birthdate",
					"address",
				}).
				AddRow(
					patient.ID,
					patient.FirstName,
					patient.LastName,
					patient.MiddleName,
					patient.Title,
					patient.BirthDate,
					patient.Address,
				)

			mock.ExpectQuery(regexp.QuoteMeta(sqlFind)).
				WithArgs(patient.ID).
				WillReturnRows(rows)

			result, err := repository.PatientById(patient.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result).Should(Equal(patient))
		})

		It("empty", func() {
			const searchTerm = "On"
			const sqlFind = `
        SELECT * FROM "patients"
        WHERE firstname LIKE ($1) OR
        WHERE lastname LIKE ($1) OR
        WHERE middlename LIKE ($1)
        `
			mock.ExpectQuery(regexp.QuoteMeta(sqlFind)).
				WithArgs(searchTerm)

			_, err := repository.FindPatients(searchTerm)
			Expect(err).Should(Equal(gorm.ErrRecordNotFound))
		})
	})

	Context("read by id", func() {
		It("found", func() {
			patient := &Patient{
				FirstName:  "One",
				LastName:   "PatientOne",
				MiddleName: "First",
				Title:      "Mr",
				BirthDate:  datatypes.Date(time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)),
				Address:    "One Address",
			}
			rows := sqlmock.
				NewRows(
					[]string{
						"firstname",
						"lastname",
						"middlename",
						"title",
						"birthdate",
						"address",
					}).
				AddRow(
					patient.ID,
					patient.FirstName,
					patient.LastName,
					patient.MiddleName,
					patient.Title,
					patient.BirthDate,
					patient.Address,
				)

			const sqlSelectOne = `
        SELECT * FROM "patients"
        WHERE (id = $1)
        ORDER BY "patients"."id"
        ASC
        LIMIT 1
        `

			mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
				WithArgs(patient.ID).
				WillReturnRows(rows)

			result, err := repository.PatientById(patient.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result).Should(Equal(patient))
		})

		It("not found", func() {
			mock.ExpectQuery(`.+`).WillReturnRows(sqlmock.NewRows(nil))
			_, err := repository.PatientById(1)
			Expect(err).Should(Equal(gorm.ErrRecordNotFound))
		})

	})

	Context("create", func() {
		var patient *Patient
		BeforeEach(func() {
			patient = &Patient{
				FirstName:  "One",
				LastName:   "PatientOne",
				MiddleName: "First",
				Title:      "Mr",
				BirthDate:  datatypes.Date(time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)),
				Address:    "One Address",
			}
		})

		It("successful create", func() {
			const sqlInsert = `
        INSERT INTO "patients" ("firstname", "lastname", "middlename", "title", "birthdate", "address")
        VALUES ($1, $2, $3, $4, $5, $6) RETURNING "patients"."id"
        `
			const newId = 1
			mock.ExpectBegin()
			mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
				WithArgs(
					patient.FirstName,
					patient.LastName,
					patient.MiddleName,
					patient.Title,
					patient.BirthDate,
					patient.Address).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newId))
			mock.ExpectCommit()

			Expect(patient.ID).Should(BeZero())

			err := repository.CreatePatient(patient)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(patient.ID).Should(BeEquivalentTo(newId))
		})

		// TODO: to write expected fail	It("expected fail create", func() {})

	})
	Context("update", func() {
		// TODO: write test
	})

	Context("delete", func() {
		// TODO: write test
	})

})
