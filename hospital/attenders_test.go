package hospital_test

import (
	"strings"
	"testing"

	"costenos-narcolepsia/hospital"
)

func newStaff(t *testing.T) (*hospital.Doctor, *hospital.Orderly) {
	t.Helper()
	d, err := hospital.NewDoctor("D-01", "Karen Ospina", 40, "Somnología")
	if err != nil {
		t.Fatal(err)
	}
	o, err := hospital.NewOrderly("O-01", "Pedro Pérez", 30, "Diurno")
	if err != nil {
		t.Fatal(err)
	}
	return d, o
}

// Polimorfismo: Doctor y Orderly viven en el mismo []Attender.
func TestAttender_Polymorphism(t *testing.T) {
	d, o := newStaff(t)
	staff := []hospital.Attender{d, o}

	p := hospital.NewPatient("P-001", "Wilfrido Berrio", 50, hospital.Severe)
	p.SufferSleepAttack("cafetería")

	for _, a := range staff {
		rec, err := a.Attend(p, p.Location())
		if err != nil {
			t.Fatalf("%s: error inesperado: %v", a.Name(), err)
		}
		if rec.AttenderID() != a.ID() {
			t.Errorf("attender esperado %s, obtenido %s", a.ID(), rec.AttenderID())
		}
	}
	if len(p.Episodes()) != 2 {
		t.Errorf("el paciente debía tener 2 episodios, tiene %d", len(p.Episodes()))
	}
}

func TestDoctor_MyEpisodesAndRoom(t *testing.T) {
	d, _ := newStaff(t)
	p := hospital.NewPatient("P-001", "Yeimy Padilla", 33, hospital.Severe)
	room := hospital.NewRoom(101, 1)
	p.SufferSleepAttack("radiología")
	if err := room.Occupy(p); err != nil {
		t.Fatal(err)
	}
	p.SetRoom(room)
	p.SetState(hospital.AsleepInBed)

	if err := d.AttendSleepEmergency(p); err != nil {
		t.Fatal(err)
	}
	eps := d.MyEpisodes()
	if len(eps) != 1 {
		t.Fatalf("se esperaba 1 episodio, hay %d", len(eps))
	}
	if !strings.Contains(eps[0].Summary(), "habitación 101") {
		t.Errorf("Summary sin habitación: %q", eps[0].Summary())
	}
	// MyEpisodes retorna copia: modificarla no altera el original.
	eps[0] = hospital.EpisodeRecord{}
	if d.MyEpisodes()[0].ID() == "" {
		t.Error("MyEpisodes debe retornar una copia defensiva")
	}
}

func TestAttend_AwakePatientFails(t *testing.T) {
	d, o := newStaff(t)
	p := hospital.NewPatient("P-002", "Despierto", 20, hospital.Mild)
	if _, err := d.Attend(p, "pasillo"); err != hospital.ErrPatientAwake {
		t.Errorf("doctor: se esperaba ErrPatientAwake, obtenido %v", err)
	}
	if _, err := o.Attend(p, "pasillo"); err != hospital.ErrPatientAwake {
		t.Errorf("camillero: se esperaba ErrPatientAwake, obtenido %v", err)
	}
	if _, err := d.Attend(nil, "pasillo"); err != hospital.ErrNilPatient {
		t.Errorf("se esperaba ErrNilPatient, obtenido %v", err)
	}
}

func TestDoctor_DiagnosePatientNoDuplicates(t *testing.T) {
	d, _ := newStaff(t)
	p := hospital.NewPatient("P-003", "Rafael", 40, hospital.Moderate)
	_ = d.DiagnosePatient(p)
	_ = d.DiagnosePatient(p)
	if n := len(d.PatientsUnderCare()); n != 1 {
		t.Errorf("se esperaba 1 paciente a cargo, hay %d", n)
	}
	if p.AssignedDoctor() == nil || p.AssignedDoctor().ID() != d.ID() {
		t.Error("el doctor asignado del paciente no se actualizó")
	}
}
