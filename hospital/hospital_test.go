package hospital_test

import (
	"strings"
	"testing"

	"costenos-narcolepsia/hospital"
)

// TestRoom_OccupySuccess valida que un paciente ingrese a una habitación disponible,
// que la habitación pase a estado Occupied al alcanzar su capacidad máxima, y que
// la lista de occupants contenga al paciente.
func TestRoom_OccupySuccess(t *testing.T) {
	room := hospital.NewRoom(101, 1)
	patient := hospital.NewPatient("P01", "Diomedes Díaz", 45, hospital.Moderate)

	if !room.IsAvailable() {
		t.Fatalf("se esperaba que la habitación %d estuviera disponible inicialmente", room.Number())
	}

	err := room.Occupy(patient)
	if err != nil {
		t.Fatalf("error inesperado al ingresar paciente a la habitación: %v", err)
	}

	if room.State() != hospital.Occupied {
		t.Errorf("estado esperado %v, obtenido %v", hospital.Occupied, room.State())
	}

	if room.IsAvailable() {
		t.Errorf("se esperaba que IsAvailable() retornara false al llenarse la habitación")
	}

	occupants := room.Occupants()
	if len(occupants) != 1 {
		t.Fatalf("se esperaba 1 ocupante, encontrados %d", len(occupants))
	}

	if occupants[0].ID() != patient.ID() {
		t.Errorf("ID del ocupante esperado %s, obtenido %s", patient.ID(), occupants[0].ID())
	}
}

// TestRoom_OccupyFullError llena una habitación de capacidad 1 con un paciente,
// intenta ingresar un segundo paciente y valida que retorne un error explicativo
// sin causar ningún pánico.
func TestRoom_OccupyFullError(t *testing.T) {
	room := hospital.NewRoom(102, 1)
	p1 := hospital.NewPatient("P01", "Poncho Zuleta", 50, hospital.Severe)
	p2 := hospital.NewPatient("P02", "Jorge Oñate", 55, hospital.Mild)

	if err := room.Occupy(p1); err != nil {
		t.Fatalf("error inesperado al alojar el primer paciente: %v", err)
	}

	err := room.Occupy(p2)
	if err == nil {
		t.Fatalf("se esperaba error al ingresar paciente en habitación sin disponibilidad, pero err fue nil")
	}

	expectedSub := "llena o no disponible"
	if !strings.Contains(err.Error(), expectedSub) {
		t.Errorf("el mensaje de error '%s' debe contener '%s'", err.Error(), expectedSub)
	}
}

// TestRoom_Release ingresa un paciente, lo libera mediante Release,
// verifica que occupants quede vacío y que el estado vuelva a ser Available.
func TestRoom_Release(t *testing.T) {
	room := hospital.NewRoom(103, 1)
	patient := hospital.NewPatient("P01", "Rafael Orozco", 38, hospital.Mild)

	if err := room.Occupy(patient); err != nil {
		t.Fatalf("error al ingresar paciente: %v", err)
	}

	err := room.Release(patient)
	if err != nil {
		t.Fatalf("error inesperado al liberar paciente: %v", err)
	}

	if len(room.Occupants()) != 0 {
		t.Errorf("se esperaba que la lista de ocupantes estuviera vacía, tiene %d", len(room.Occupants()))
	}

	if room.State() != hospital.Available {
		t.Errorf("estado esperado %v, obtenido %v", hospital.Available, room.State())
	}

	if !room.IsAvailable() {
		t.Errorf("IsAvailable() debe retornar true tras la liberación")
	}
}

// Tests complementarios de Persona 1: Person y Patient
func TestPatient_SufferSleepAttackAndWakeUp(t *testing.T) {
	patient := hospital.NewPatient("P03", "Silvestre Dangond", 42, hospital.Severe)

	if patient.State() != hospital.Awake {
		t.Fatalf("el paciente debe inicializarse con estado Awake")
	}

	patient.SufferSleepAttack("Pasillo de Urgencias")
	if patient.State() != hospital.AsleepInHallway {
		t.Errorf("estado esperado AsleepInHallway, obtenido %v", patient.State())
	}
	if patient.Location() != "Pasillo de Urgencias" {
		t.Errorf("ubicación esperada 'Pasillo de Urgencias', obtenida '%s'", patient.Location())
	}

	// Despertar al paciente
	if err := patient.WakeUp(); err != nil {
		t.Fatalf("error inesperado al despertar al paciente: %v", err)
	}
	if patient.State() != hospital.Awake {
		t.Errorf("estado esperado Awake, obtenido %v", patient.State())
	}
	if patient.Location() != "" {
		t.Errorf("la ubicación debía limpiarse, contiene '%s'", patient.Location())
	}

	// Intentar despertar a un paciente ya despierto debe retornar error
	err := patient.WakeUp()
	if err == nil {
		t.Fatalf("se esperaba error al despertar a un paciente que ya está despierto")
	}
}

func TestPerson_Getters(t *testing.T) {
	person := hospital.NewPerson("PERS-01", "Kaleth Morales", 22)
	if person.ID() != "PERS-01" {
		t.Errorf("ID esperado 'PERS-01', obtenido '%s'", person.ID())
	}
	if person.Name() != "Kaleth Morales" {
		t.Errorf("Nombre esperado 'Kaleth Morales', obtenido '%s'", person.Name())
	}
	if person.Age() != 22 {
		t.Errorf("Edad esperada 22, obtenida %d", person.Age())
	}
}
