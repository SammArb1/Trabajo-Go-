package hospital

import "errors"

// Constantes y tipos enumerados para el dominio hospitalario.

// PatientState representa los posibles estados de vigilia y sueño de un paciente.
type PatientState int

const (
	Awake PatientState = iota
	AsleepInHallway
	AsleepInBed
)

// String implementa fmt.Stringer para PatientState.
func (s PatientState) String() string {
	switch s {
	case Awake:
		return "Awake"
	case AsleepInHallway:
		return "AsleepInHallway"
	case AsleepInBed:
		return "AsleepInBed"
	default:
		return "Unknown"
	}
}

// NarcolepsyLevel representa el grado de severidad del trastorno de narcolepsia.
type NarcolepsyLevel int

const (
	Mild NarcolepsyLevel = iota
	Moderate
	Severe
)

// String implementa fmt.Stringer para NarcolepsyLevel.
func (l NarcolepsyLevel) String() string {
	switch l {
	case Mild:
		return "Mild"
	case Moderate:
		return "Moderate"
	case Severe:
		return "Severe"
	default:
		return "Unknown"
	}
}

// RoomState modela la disponibilidad de una habitación en el hospital.
type RoomState int

const (
	Available RoomState = iota
	Occupied
)

// String implementa fmt.Stringer para RoomState.
func (r RoomState) String() string {
	switch r {
	case Available:
		return "Available"
	case Occupied:
		return "Occupied"
	default:
		return "Unknown"
	}
}

// Errores de dominio para control defensivo sin pánico.
var (
	ErrEmptyID          = errors.New("el identificador no puede estar vacío")
	ErrEmptyName        = errors.New("el nombre no puede estar vacío")
	ErrInvalidAge       = errors.New("la edad no puede ser negativa")
	ErrNilPatient       = errors.New("el paciente no puede ser nil")
	ErrNilRoom          = errors.New("la habitación no puede ser nil")
	ErrRoomOccupied     = errors.New("la habitación ya está ocupada")
	ErrRoomNotOccupied  = errors.New("la habitación no está ocupada")
	ErrPatientNotFound  = errors.New("paciente no encontrado")
	ErrRoomNotFound     = errors.New("habitación no encontrada")
	ErrAttenderNotFound = errors.New("atendedor o personal no encontrado")
)
