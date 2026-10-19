package hospital

import (
	"fmt"
	"sync/atomic"
	"time"
)

// episodeCounter genera identificadores secuenciales y únicos de episodio.
var episodeCounter atomic.Int64

// EpisodeRecord almacena la información inmutable de un episodio de narcolepsia atendido.
// Es un tipo valor: sus métodos tienen receptor por valor y no pueden modificar el original.
type EpisodeRecord struct {
	id         string
	patientID  string
	attenderID string
	location   string
	room       string // vacío si el paciente no tuvo habitación
	timestamp  time.Time
	notes      string
	resolved   bool
}

// NewEpisodeRecord crea una nueva instancia de EpisodeRecord con campos no exportados.
func NewEpisodeRecord(id, patientID, attenderID, location string, timestamp time.Time, notes string, resolved bool) (EpisodeRecord, error) {
	if id == "" || patientID == "" || attenderID == "" {
		return EpisodeRecord{}, ErrEmptyID
	}
	return EpisodeRecord{
		id:         id,
		patientID:  patientID,
		attenderID: attenderID,
		location:   location,
		timestamp:  timestamp,
		notes:      notes,
		resolved:   resolved,
	}, nil
}

// newAttendRecord construye el registro de una atención y lo agrega al historial
// del paciente. Lo comparten Doctor y Orderly para no duplicar lógica.
func newAttendRecord(a Attender, p *Patient, location, notes string) (EpisodeRecord, error) {
	if p == nil {
		return EpisodeRecord{}, ErrNilPatient
	}
	if p.State() == Awake {
		return EpisodeRecord{}, ErrPatientAwake
	}
	id := fmt.Sprintf("EP-%04d", episodeCounter.Add(1))
	rec, err := NewEpisodeRecord(id, p.ID(), a.ID(), location, time.Now(), notes, true)
	if err != nil {
		return EpisodeRecord{}, err
	}
	if r := p.Room(); r != nil {
		rec = rec.WithRoom(r.ID())
	}
	if err := p.AddEpisode(rec); err != nil {
		return EpisodeRecord{}, err
	}
	return rec, nil
}

// WithRoom retorna una copia del registro con la habitación asignada.
func (e EpisodeRecord) WithRoom(room string) EpisodeRecord {
	e.room = room
	return e
}

// Summary retorna una línea legible: fecha, paciente, lugar y habitación.
func (e EpisodeRecord) Summary() string {
	room := "sin habitación"
	if e.room != "" {
		room = "habitación " + e.room
	}
	return fmt.Sprintf("%s %s %s -> %s",
		e.timestamp.Format("2006-01-02 15:04"), e.patientID, e.location, room)
}

// ID retorna el identificador único del episodio.
func (e EpisodeRecord) ID() string { return e.id }

// PatientID retorna el identificador del paciente atendido.
func (e EpisodeRecord) PatientID() string { return e.patientID }

// AttenderID retorna el identificador del profesional o camillero que prestó auxilio.
func (e EpisodeRecord) AttenderID() string { return e.attenderID }

// Location retorna el lugar donde el paciente se quedó dormido.
func (e EpisodeRecord) Location() string { return e.location }

// Room retorna el número de habitación asignada ("" si no hubo).
func (e EpisodeRecord) Room() string { return e.room }

// Timestamp retorna la marca temporal de cuando ocurrió la atención.
func (e EpisodeRecord) Timestamp() time.Time { return e.timestamp }

// Notes retorna observaciones clínicas o de traslado registradas.
func (e EpisodeRecord) Notes() string { return e.notes }

// IsResolved indica si el episodio fue mitigado exitosamente.
func (e EpisodeRecord) IsResolved() bool { return e.resolved }
