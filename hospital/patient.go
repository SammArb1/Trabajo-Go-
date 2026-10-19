package hospital

import (
	"fmt"
)

// Patient representa a un paciente que padece narcolepsia en el hospital.
// Embebe a Person mediante composición y protege su estado interno con campos privados.
type Patient struct {
	Person
	level           NarcolepsyLevel
	state           PatientState
	currentLocation string
	assignedDoctor  Attender
	room            *Room
	episodes        []EpisodeRecord
}

// NewPatient construye un nuevo paciente inicializando su estado en Awake.
func NewPatient(id, name string, age int, level NarcolepsyLevel) *Patient {
	return &Patient{
		Person:          NewPerson(id, name, age),
		level:           level,
		state:           Awake,
		currentLocation: "",
		assignedDoctor:  nil,
		episodes:        make([]EpisodeRecord, 0),
	}
}

// SufferSleepAttack simula un ataque súbito de narcolepsia, pasando el estado a
// AsleepInHallway y guardando la ubicación física del episodio.
func (p *Patient) SufferSleepAttack(location string) {
	if p == nil {
		return
	}
	p.state = AsleepInHallway
	p.currentLocation = location
}

// WakeUp despierta al paciente restableciendo su estado a Awake y limpiando la ubicación.
// Retorna error si el paciente ya se encontraba despierto.
func (p *Patient) WakeUp() error {
	if p == nil {
		return fmt.Errorf("el paciente no puede ser nil")
	}
	if p.state == Awake {
		return fmt.Errorf("el paciente %s ya se encuentra despierto", p.Name())
	}
	p.state = Awake
	p.currentLocation = ""
	return nil
}

// State retorna el estado actual de sueño/vigilia del paciente.
func (p *Patient) State() PatientState {
	if p == nil {
		return Awake
	}
	return p.state
}

// Level retorna el nivel de severidad de narcolepsia diagnosticado.
func (p *Patient) Level() NarcolepsyLevel {
	if p == nil {
		return Mild
	}
	return p.level
}

// Location retorna la ubicación física actual del paciente.
func (p *Patient) Location() string {
	if p == nil {
		return ""
	}
	return p.currentLocation
}

// CurrentLocation retorna la ubicación física actual del paciente (alias de Location).
func (p *Patient) CurrentLocation() string {
	return p.Location()
}

// SetState permite actualizar de forma controlada el estado del paciente (ej. AsleepInBed).
func (p *Patient) SetState(s PatientState) {
	if p == nil {
		return
	}
	p.state = s
}

// SetCurrentLocation actualiza la ubicación física actual del paciente.
func (p *Patient) SetCurrentLocation(loc string) {
	if p == nil {
		return
	}
	p.currentLocation = loc
}

// AssignedDoctor retorna el personal asignado a la atención del paciente.
func (p *Patient) AssignedDoctor() Attender {
	if p == nil {
		return nil
	}
	return p.assignedDoctor
}

// SetAssignedDoctor asigna un médico o asistente al paciente.
func (p *Patient) SetAssignedDoctor(d Attender) {
	if p == nil {
		return
	}
	p.assignedDoctor = d
}

// Episodes retorna una copia de la lista de episodios registrados para el paciente.
func (p *Patient) Episodes() []EpisodeRecord {
	if p == nil {
		return nil
	}
	records := make([]EpisodeRecord, len(p.episodes))
	copy(records, p.episodes)
	return records
}

// AddEpisode registra un nuevo episodio atendido en el historial del paciente.
func (p *Patient) AddEpisode(record EpisodeRecord) error {
	if p == nil {
		return ErrNilPatient
	}
	p.episodes = append(p.episodes, record)
	return nil
}

// Room retorna la habitación donde está alojado el paciente (nil si no tiene).
func (p *Patient) Room() *Room {
	if p == nil {
		return nil
	}
	return p.room
}

// SetRoom registra la habitación asignada al paciente (nil para limpiarla).
func (p *Patient) SetRoom(r *Room) {
	if p == nil {
		return
	}
	p.room = r
}
