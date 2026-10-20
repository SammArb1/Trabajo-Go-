package hospital_test

import (
	"errors"
	"testing"

	"costenos-narcolepsia/hospital"
)

// newScenario crea un hospital con 1 médico, 1 camillero, las habitaciones pedidas
// y 3 pacientes (P-002 es Severe).
func newScenario(t *testing.T, rooms int) (*hospital.Hospital, []*hospital.Patient) {
	t.Helper()
	h := hospital.NewHospital("Test")
	d, err := hospital.NewDoctor("D-01", "Karen", 40, "Neuro")
	if err != nil {
		t.Fatal(err)
	}
	o, err := hospital.NewOrderly("O-01", "Jairo", 30, "Diurno")
	if err != nil {
		t.Fatal(err)
	}
	if err := h.HireDoctor(d); err != nil {
		t.Fatal(err)
	}
	if err := h.HireOrderly(o); err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= rooms; i++ {
		if err := h.AddRoom(hospital.NewRoom(100+i, 1)); err != nil {
			t.Fatal(err)
		}
	}
	ps := []*hospital.Patient{
		hospital.NewPatient("P-001", "Uno", 30, hospital.Mild),
		hospital.NewPatient("P-002", "Dos", 30, hospital.Severe),
		hospital.NewPatient("P-003", "Tres", 30, hospital.Moderate),
	}
	for _, p := range ps {
		if err := h.AdmitPatient(p); err != nil {
			t.Fatal(err)
		}
	}
	return h, ps
}

func TestAssignRoom_FreeRoom(t *testing.T) {
	h, ps := newScenario(t, 2)
	ps[0].SufferSleepAttack("cafeteria")
	room, err := h.AssignRoom(ps[0])
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if room.Number() != 101 {
		t.Errorf("debía asignar la habitación 101, asignó %d", room.Number())
	}
	if room.State() != hospital.Occupied || ps[0].State() != hospital.AsleepInBed {
		t.Errorf("room=%v patient=%v", room.State(), ps[0].State())
	}
}

func TestAssignRoom_NoRoomAvailable(t *testing.T) {
	h, ps := newScenario(t, 0)
	ps[0].SufferSleepAttack("pasillo")
	room, err := h.AssignRoom(ps[0])
	if !errors.Is(err, hospital.ErrNoRoomAvailable) || room != nil {
		t.Fatalf("room=%v err=%v", room, err)
	}
	if ps[0].State() != hospital.AsleepInHallway {
		t.Errorf("el paciente debía seguir en el pasillo, está %v", ps[0].State())
	}
}

func TestPatientsInHallway(t *testing.T) {
	h, ps := newScenario(t, 1)
	if _, err := h.HandleSleepAttack(ps[0], "a"); err != nil {
		t.Fatal(err)
	}
	if _, err := h.HandleSleepAttack(ps[1], "b"); !errors.Is(err, hospital.ErrNoRoomAvailable) {
		t.Fatalf("se esperaba ErrNoRoomAvailable, got %v", err)
	}
	got := h.PatientsInHallway()
	if len(got) != 1 || got[0] != ps[1] {
		t.Fatalf("solo P-002 debía estar en el pasillo, got %v", got)
	}
}

func TestWakePatientFreesRoomForNext(t *testing.T) {
	h, ps := newScenario(t, 1)
	_, _ = h.HandleSleepAttack(ps[0], "a")
	_, _ = h.HandleSleepAttack(ps[1], "b")
	if err := h.WakePatient(ps[0]); err != nil {
		t.Fatal(err)
	}
	room, err := h.AssignRoom(ps[1])
	if err != nil || room.Number() != 101 {
		t.Fatalf("room=%v err=%v", room, err)
	}
}

func TestSevereReportAndDispatch(t *testing.T) {
	h, ps := newScenario(t, 3)
	for _, p := range ps {
		if _, err := h.HandleSleepAttack(p, "x"); err != nil {
			t.Fatal(err)
		}
	}
	report := h.SevereReport()
	if len(report) != 1 || report[ps[1]] != 1 {
		t.Errorf("reporte inesperado: %v", report)
	}
	// Orden de despacho: médico, camillero, médico.
	if n := len(h.EpisodesByAttender("D-01")); n != 2 {
		t.Errorf("episodios de D-01 = %d, se esperaban 2", n)
	}
	if n := len(h.EpisodesByAttender("O-01")); n != 1 {
		t.Errorf("episodios de O-01 = %d, se esperaba 1", n)
	}
}
