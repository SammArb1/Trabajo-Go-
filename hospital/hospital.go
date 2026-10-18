package hospital

import (
	"errors"
)

// Hospital es el agregador principal del dominio.
// Encapsula todas las colecciones y entidades del "Hospital de los Costeños con Narcolepsia".
type Hospital struct {
	name      string
	rooms     map[string]*Room
	patients  map[string]*Patient
	doctors   map[string]*Doctor
	orderlies map[string]*Orderly
	episodes  []EpisodeRecord
}

// NewHospital inicializa una nueva instancia del hospital con mapas y listas vacíos.
func NewHospital(name string) *Hospital {
	if name == "" {
		name = "Hospital de los Costeños con Narcolepsia"
	}
	return &Hospital{
		name:      name,
		rooms:     make(map[string]*Room),
		patients:  make(map[string]*Patient),
		doctors:   make(map[string]*Doctor),
		orderlies: make(map[string]*Orderly),
		episodes:  make([]EpisodeRecord, 0),
	}
}

// Name retorna el nombre de la institución hospitalaria.
func (h *Hospital) Name() string {
	if h == nil {
		return ""
	}
	return h.name
}

// SetName actualiza el nombre del hospital validando que no sea vacío.
func (h *Hospital) SetName(name string) error {
	if h == nil {
		return errors.New("hospital no puede ser nil")
	}
	if name == "" {
		return ErrEmptyName
	}
	h.name = name
	return nil
}

// AddRoom registra una nueva habitación en el hospital.
func (h *Hospital) AddRoom(room *Room) error {
	if h == nil {
		return errors.New("hospital no puede ser nil")
	}
	if room == nil {
		return ErrNilRoom
	}
	if _, exists := h.rooms[room.ID()]; exists {
		return errors.New("la habitación ya existe en el hospital")
	}
	h.rooms[room.ID()] = room
	return nil
}

// Room busca una habitación por su identificador.
func (h *Hospital) Room(id string) (*Room, error) {
	if h == nil {
		return nil, errors.New("hospital no puede ser nil")
	}
	r, exists := h.rooms[id]
	if !exists {
		return nil, ErrRoomNotFound
	}
	return r, nil
}

// Rooms retorna un slice con todas las habitaciones registradas.
func (h *Hospital) Rooms() []*Room {
	if h == nil {
		return nil
	}
	list := make([]*Room, 0, len(h.rooms))
	for _, r := range h.rooms {
		list = append(list, r)
	}
	return list
}

// RegisterPatient añade un nuevo paciente con narcolepsia al registro general.
func (h *Hospital) RegisterPatient(patient *Patient) error {
	if h == nil {
		return errors.New("hospital no puede ser nil")
	}
	if patient == nil {
		return ErrNilPatient
	}
	if _, exists := h.patients[patient.ID()]; exists {
		return errors.New("el paciente ya se encuentra registrado")
	}
	h.patients[patient.ID()] = patient
	return nil
}

// Patient busca un paciente por su identificador.
func (h *Hospital) Patient(id string) (*Patient, error) {
	if h == nil {
		return nil, errors.New("hospital no puede ser nil")
	}
	p, exists := h.patients[id]
	if !exists {
		return nil, ErrPatientNotFound
	}
	return p, nil
}

// Patients retorna un slice con todos los pacientes registrados.
func (h *Hospital) Patients() []*Patient {
	if h == nil {
		return nil
	}
	list := make([]*Patient, 0, len(h.patients))
	for _, p := range h.patients {
		list = append(list, p)
	}
	return list
}

// RegisterDoctor registra un médico especialista en el hospital.
func (h *Hospital) RegisterDoctor(doctor *Doctor) error {
	if h == nil {
		return errors.New("hospital no puede ser nil")
	}
	if doctor == nil {
		return ErrAttenderNotFound
	}
	if _, exists := h.doctors[doctor.ID()]; exists {
		return errors.New("el médico ya se encuentra registrado")
	}
	h.doctors[doctor.ID()] = doctor
	return nil
}

// Doctor busca un médico por su identificador.
func (h *Hospital) Doctor(id string) (*Doctor, error) {
	if h == nil {
		return nil, errors.New("hospital no puede ser nil")
	}
	d, exists := h.doctors[id]
	if !exists {
		return nil, ErrAttenderNotFound
	}
	return d, nil
}

// Doctors retorna la lista de médicos registrados.
func (h *Hospital) Doctors() []*Doctor {
	if h == nil {
		return nil
	}
	list := make([]*Doctor, 0, len(h.doctors))
	for _, d := range h.doctors {
		list = append(list, d)
	}
	return list
}

// RegisterOrderly registra un camillero en el hospital.
func (h *Hospital) RegisterOrderly(orderly *Orderly) error {
	if h == nil {
		return errors.New("hospital no puede ser nil")
	}
	if orderly == nil {
		return ErrAttenderNotFound
	}
	if _, exists := h.orderlies[orderly.ID()]; exists {
		return errors.New("el camillero ya se encuentra registrado")
	}
	h.orderlies[orderly.ID()] = orderly
	return nil
}

// Orderly busca un camillero por su identificador.
func (h *Hospital) Orderly(id string) (*Orderly, error) {
	if h == nil {
		return nil, errors.New("hospital no puede ser nil")
	}
	o, exists := h.orderlies[id]
	if !exists {
		return nil, ErrAttenderNotFound
	}
	return o, nil
}

// Orderlies retorna la lista de camilleros registrados.
func (h *Hospital) Orderlies() []*Orderly {
	if h == nil {
		return nil
	}
	list := make([]*Orderly, 0, len(h.orderlies))
	for _, o := range h.orderlies {
		list = append(list, o)
	}
	return list
}

// Episodes retorna el historial de todos los episodios registrados en el hospital.
func (h *Hospital) Episodes() []EpisodeRecord {
	if h == nil {
		return nil
	}
	list := make([]EpisodeRecord, len(h.episodes))
	copy(list, h.episodes)
	return list
}

// RecordEpisode almacena un nuevo registro de episodio en el hospital.
func (h *Hospital) RecordEpisode(record EpisodeRecord) error {
	if h == nil {
		return errors.New("hospital no puede ser nil")
	}
	h.episodes = append(h.episodes, record)
	return nil
}
