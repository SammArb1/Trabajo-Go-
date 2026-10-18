package hospital

import (
	"fmt"
	"time"
)

// Doctor representa a un médico especialista del hospital.
// Embebe a Person e implementa la interfaz Attender.
type Doctor struct {
	Person
	specialty string
}

// Verificación en tiempo de compilación de que *Doctor satisface la interfaz Attender.
var _ Attender = (*Doctor)(nil)

// NewDoctor crea un nuevo médico validando los datos requeridos.
func NewDoctor(id, name string, age int, specialty string) (*Doctor, error) {
	person, err := NewPerson(id, name, age)
	if err != nil {
		return nil, err
	}
	if specialty == "" {
		specialty = "Medicina General / Somnología"
	}
	return &Doctor{
		Person:    *person,
		specialty: specialty,
	}, nil
}

// Specialty retorna la especialidad médica del doctor.
func (d *Doctor) Specialty() string {
	if d == nil {
		return ""
	}
	return d.specialty
}

// SetSpecialty actualiza la especialidad médica previa validación.
func (d *Doctor) SetSpecialty(specialty string) error {
	if d == nil {
		return ErrAttenderNotFound
	}
	if specialty == "" {
		return ErrEmptyName
	}
	d.specialty = specialty
	return nil
}

// Attend atiende medicamente a un paciente en episodio de narcolepsia.
// Cumple con la interfaz Attender.
func (d *Doctor) Attend(p *Patient, location string) (EpisodeRecord, error) {
	if d == nil {
		return EpisodeRecord{}, ErrAttenderNotFound
	}
	if p == nil {
		return EpisodeRecord{}, ErrNilPatient
	}

	notes := fmt.Sprintf("Atención médica por Dr. %s (%s). Paciente valorado y estabilizado.", d.Name(), d.specialty)
	recordID := fmt.Sprintf("EP-DOC-%s-%d", d.ID(), time.Now().UnixNano())

	record, err := NewEpisodeRecord(recordID, p.ID(), d.ID(), location, time.Now(), notes, true)
	if err != nil {
		return EpisodeRecord{}, err
	}

	if err := p.AddEpisode(record); err != nil {
		return EpisodeRecord{}, err
	}

	return record, nil
}
