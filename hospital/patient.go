package hospital

// Patient representa a un paciente que padece narcolepsia en el hospital.
// Embebe a Person para reutilizar la identidad y validaciones básicas.
type Patient struct {
	Person
	state    PatientState
	level    NarcolepsyLevel
	roomID   string
	episodes []EpisodeRecord
}

// NewPatient construye un nuevo paciente asegurando datos válidos.
func NewPatient(id, name string, age int, level NarcolepsyLevel) (*Patient, error) {
	person, err := NewPerson(id, name, age)
	if err != nil {
		return nil, err
	}
	return &Patient{
		Person:   *person,
		state:    Awake,
		level:    level,
		roomID:   "",
		episodes: make([]EpisodeRecord, 0),
	}, nil
}

// State retorna el estado actual de sueño/vigilia del paciente.
func (p *Patient) State() PatientState {
	if p == nil {
		return Awake
	}
	return p.state
}

// SetState actualiza el estado del paciente.
func (p *Patient) SetState(s PatientState) error {
	if p == nil {
		return ErrNilPatient
	}
	p.state = s
	return nil
}

// Level retorna la gravedad de narcolepsia diagnosticada.
func (p *Patient) Level() NarcolepsyLevel {
	if p == nil {
		return Mild
	}
	return p.level
}

// SetLevel modifica la gravedad de narcolepsia diagnosticada.
func (p *Patient) SetLevel(l NarcolepsyLevel) error {
	if p == nil {
		return ErrNilPatient
	}
	p.level = l
	return nil
}

// RoomID retorna el código de la habitación donde reposa el paciente, si aplica.
func (p *Patient) RoomID() string {
	if p == nil {
		return ""
	}
	return p.roomID
}

// SetRoomID actualiza la habitación vinculada al paciente.
func (p *Patient) SetRoomID(roomID string) error {
	if p == nil {
		return ErrNilPatient
	}
	p.roomID = roomID
	return nil
}

// Episodes devuelve una copia defensiva del histórico de episodios del paciente.
func (p *Patient) Episodes() []EpisodeRecord {
	if p == nil {
		return nil
	}
	records := make([]EpisodeRecord, len(p.episodes))
	copy(records, p.episodes)
	return records
}

// AddEpisode agrega un nuevo registro de atención al historial del paciente.
func (p *Patient) AddEpisode(record EpisodeRecord) error {
	if p == nil {
		return ErrNilPatient
	}
	p.episodes = append(p.episodes, record)
	return nil
}

// FallAsleep cambia el estado del paciente a dormido según el lugar donde se encuentre.
func (p *Patient) FallAsleep(location string) error {
	if p == nil {
		return ErrNilPatient
	}
	if p.roomID != "" && location == p.roomID {
		p.state = AsleepInBed
	} else {
		p.state = AsleepInHallway
	}
	return nil
}

// WakeUp despierta al paciente restableciendo su estado a Awake.
func (p *Patient) WakeUp() error {
	if p == nil {
		return ErrNilPatient
	}
	p.state = Awake
	return nil
}
