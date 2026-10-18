package hospital

import (
	"errors"
	"fmt"
	"strconv"
)

// Room modela una habitación del hospital con capacidad definida y control de ocupación.
// Todos sus campos son privados para garantizar encapsulamiento riguroso.
type Room struct {
	number    int
	capacity  int
	state     RoomState
	occupants []*Patient
}

// NewRoom construye una nueva habitación con número y capacidad, inicializada como Available.
func NewRoom(number, capacity int) *Room {
	return &Room{
		number:    number,
		capacity:  capacity,
		state:     Available,
		occupants: make([]*Patient, 0, capacity),
	}
}

// IsAvailable retorna true si la habitación cuenta con cupo y su estado es Available.
func (r *Room) IsAvailable() bool {
	if r == nil {
		return false
	}
	return len(r.occupants) < r.capacity && r.state == Available
}

// Occupy ingresa un paciente a la habitación si hay disponibilidad.
// Retorna error descriptivo si el paciente es nil o si la habitación está llena o no disponible.
// Si se llena la capacidad máxima, actualiza su estado a Occupied.
func (r *Room) Occupy(p *Patient) error {
	if r == nil {
		return errors.New("la habitación no puede ser nil")
	}
	if p == nil {
		return errors.New("el paciente no puede ser nil")
	}
	if !r.IsAvailable() {
		return fmt.Errorf("habitacion %d llena o no disponible", r.number)
	}

	r.occupants = append(r.occupants, p)
	if len(r.occupants) == r.capacity {
		r.state = Occupied
	}
	return nil
}

// Release busca y retira a un paciente de la habitación.
// Si no lo encuentra, retorna error. Si tras la salida hay cupo, pasa su estado a Available.
func (r *Room) Release(p *Patient) error {
	if r == nil {
		return errors.New("la habitación no puede ser nil")
	}
	if p == nil {
		return errors.New("el paciente no puede ser nil")
	}

	index := -1
	for i, occ := range r.occupants {
		if occ == p || (occ != nil && occ.ID() == p.ID()) {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("el paciente no se encuentra en la habitacion %d", r.number)
	}

	r.occupants = append(r.occupants[:index], r.occupants[index+1:]...)
	if len(r.occupants) < r.capacity {
		r.state = Available
	}
	return nil
}

// Number retorna el número identificador de la habitación.
func (r *Room) Number() int {
	if r == nil {
		return 0
	}
	return r.number
}

// Capacity retorna la capacidad máxima de pacientes en la habitación.
func (r *Room) Capacity() int {
	if r == nil {
		return 0
	}
	return r.capacity
}

// State retorna el estado actual de la habitación (Available u Occupied).
func (r *Room) State() RoomState {
	if r == nil {
		return Available
	}
	return r.state
}

// Occupants retorna una copia defensiva de la lista de pacientes actualmente alojados.
func (r *Room) Occupants() []*Patient {
	if r == nil {
		return nil
	}
	result := make([]*Patient, len(r.occupants))
	copy(result, r.occupants)
	return result
}

// ID retorna el identificador de la habitación en formato string para compatibilidad.
func (r *Room) ID() string {
	if r == nil {
		return ""
	}
	return strconv.Itoa(r.number)
}
