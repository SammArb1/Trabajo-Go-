package hospital

import "fmt"

// Orderly representa a un camillero: socorre y traslada pacientes dormidos,
// pero no diagnostica. Embebe a Person e implementa Attender.
type Orderly struct {
	Person
	shift    string
	episodes []EpisodeRecord
}

// Verificación en tiempo de compilación de que *Orderly satisface Attender.
var _ Attender = (*Orderly)(nil)

// NewOrderly construye una nueva instancia de camillero validando sus datos.
func NewOrderly(id, name string, age int, shift string) (*Orderly, error) {
	if id == "" {
		return nil, ErrEmptyID
	}
	if name == "" {
		return nil, ErrEmptyName
	}
	if age < 0 {
		return nil, ErrInvalidAge
	}
	if shift == "" {
		shift = "Diurno"
	}
	return &Orderly{
		Person: NewPerson(id, name, age),
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

// Attend asiste y moviliza a un paciente dormido sin diagnosticarlo.
// Cumple la interfaz Attender.
func (o *Orderly) Attend(p *Patient, location string) (EpisodeRecord, error) {
	if o == nil {
		return EpisodeRecord{}, ErrAttenderNotFound
	}
	notes := fmt.Sprintf("Traslado por camillero %s (turno %s).", o.Name(), o.shift)
	rec, err := newAttendRecord(o, p, location, notes)
	if err != nil {
		return EpisodeRecord{}, err
	}
	o.episodes = append(o.episodes, rec)
	return rec, nil
}

// MyEpisodes retorna una copia de los traslados realizados por el camillero.
func (o *Orderly) MyEpisodes() []EpisodeRecord {
	if o == nil {
		return nil
	}
	out := make([]EpisodeRecord, len(o.episodes))
	copy(out, o.episodes)
	return out
}