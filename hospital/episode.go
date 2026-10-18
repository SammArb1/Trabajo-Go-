package hospital

import (
	"time"
)

// EpisodeRecord almacena la información inmutable de un episodio de narcolepsia atendido.
type EpisodeRecord struct {
	id         string
	patientID  string
	attenderID string
	location   string
	timestamp  time.Time
	notes      string
	resolved   bool
}

// NewEpisodeRecord crea una nueva instancia de EpisodeRecord con campos no exportados.
func NewEpisodeRecord(id, patientID, attenderID, location string, timestamp time.Time, notes string, resolved bool) (EpisodeRecord, error) {
	if id == "" {
		return EpisodeRecord{}, ErrEmptyID
	}
	if patientID == "" {
		return EpisodeRecord{}, ErrEmptyID
	}
	if attenderID == "" {
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

// ID retorna el identificador único del episodio.
func (e EpisodeRecord) ID() string {
	return e.id
}

// PatientID retorna el identificador del paciente atendido.
func (e EpisodeRecord) PatientID() string {
	return e.patientID
}

// AttenderID retorna el identificador del profesional o camillero que prestó auxilio.
func (e EpisodeRecord) AttenderID() string {
	return e.attenderID
}

// Location retorna la ubicación física del suceso (ej. "Pasillo Principal", "Habitación 101").
func (e EpisodeRecord) Location() string {
	return e.location
}

// Timestamp retorna la marca temporal de cuando ocurrió la atención.
func (e EpisodeRecord) Timestamp() time.Time {
	return e.timestamp
}

// Notes retorna observaciones clínicas o de traslado registradas.
func (e EpisodeRecord) Notes() string {
	return e.notes
}

// IsResolved indica si el episodio fue mitigado exitosamente.
func (e EpisodeRecord) IsResolved() bool {
	return e.resolved
}
