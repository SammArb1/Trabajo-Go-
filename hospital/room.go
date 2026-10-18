package hospital

// Room representa una habitación en el hospital con control de disponibilidad.
type Room struct {
	id               string
	state            RoomState
	currentPatientID string
}

// NewRoom construye una nueva habitación disponible.
func NewRoom(id string) (*Room, error) {
	if id == "" {
		return nil, ErrEmptyID
	}
	return &Room{
		id:               id,
		state:            Available,
		currentPatientID: "",
	}, nil
}

// ID retorna el identificador único o número de la habitación.
func (r *Room) ID() string {
	if r == nil {
		return ""
	}
	return r.id
}

// State retorna el estado actual de ocupación (Available u Occupied).
func (r *Room) State() RoomState {
	if r == nil {
		return Available
	}
	return r.state
}

// CurrentPatientID retorna el identificador del paciente asignado, o vacío si está libre.
func (r *Room) CurrentPatientID() string {
	if r == nil {
		return ""
	}
	return r.currentPatientID
}

// AssignPatient aloja a un paciente en la habitación si está disponible.
func (r *Room) AssignPatient(patientID string) error {
	if r == nil {
		return ErrNilRoom
	}
	if patientID == "" {
		return ErrEmptyID
	}
	if r.state == Occupied {
		return ErrRoomOccupied
	}
	r.currentPatientID = patientID
	r.state = Occupied
	return nil
}

// Release libera la habitación marcándola como disponible.
func (r *Room) Release() error {
	if r == nil {
		return ErrNilRoom
	}
	if r.state == Available {
		return ErrRoomNotOccupied
	}
	r.currentPatientID = ""
	r.state = Available
	return nil
}
