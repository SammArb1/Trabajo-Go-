package hospital

import (
	"fmt"
	"time"
)

// Orderly representa a un camillero responsable de socorrer y trasladar
// a los pacientes narcolépticos que caen dormidos en los pasillos.
// Embebe a Person e implementa la interfaz Attender.
type Orderly struct {
	Person
	shift string
}

// Verificación en tiempo de compilación de que *Orderly satisface la interfaz Attender.
var _ Attender = (*Orderly)(nil)

// NewOrderly construye una nueva instancia de camillero validando sus datos.
func NewOrderly(id, name string, age int, shift string) (*Orderly, error) {
	person, err := NewPerson(id, name, age)
	if err != nil {
		return nil, err
	}
	if shift == "" {
		shift = "Diurno"
	}
	return &Orderly{
		Person: *person,
		shift:  shift,
	}, nil
}

// Shift retorna el turno de trabajo asignado al camillero.
func (o *Orderly) Shift() string {
	if o == nil {
		return ""
	}
	return o.shift
}

// SetShift actualiza el turno de trabajo previa validación.
func (o *Orderly) SetShift(shift string) error {
	if o == nil {
		return ErrAttenderNotFound
	}
	if shift == "" {
		return ErrEmptyName
	}
	o.shift = shift
	return nil
}

// Attend asiste y moviliza a un paciente hacia un lugar seguro.
// Cumple con la interfaz Attender.
func (o *Orderly) Attend(p *Patient, location string) (EpisodeRecord, error) {
	if o == nil {
		return EpisodeRecord{}, ErrAttenderNotFound
	}
	if p == nil {
		return EpisodeRecord{}, ErrNilPatient
	}

	notes := fmt.Sprintf("Asistencia y movilización por camillero %s (Turno: %s).", o.Name(), o.shift)
	recordID := fmt.Sprintf("EP-ORD-%s-%d", o.ID(), time.Now().UnixNano())

	record, err := NewEpisodeRecord(recordID, p.ID(), o.ID(), location, time.Now(), notes, true)
	if err != nil {
		return EpisodeRecord{}, err
	}

	if err := p.AddEpisode(record); err != nil {
		return EpisodeRecord{}, err
	}

	return record, nil
}
