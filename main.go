package main

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"costenos-narcolepsia/hospital"
)

func main() {
	h := hospital.NewHospital("Hospital de los Costeños con Narcolepsia")

	karen := must(hospital.NewDoctor("D-01", "Karen Ospina", 41, "Neurología"))
	andres := must(hospital.NewDoctor("D-02", "Andrés Mendoza", 38, "Medicina del sueño"))
	jairo := must(hospital.NewOrderly("O-01", "Jairo Camacho", 33, "Diurno"))
	check(h.HireDoctor(karen))
	check(h.HireDoctor(andres))
	check(h.HireOrderly(jairo))

	for _, n := range []int{101, 102, 103} {
		check(h.AddRoom(hospital.NewRoom(n, 1)))
	}

	p1 := hospital.NewPatient("P-001", "Luz Marina Caro", 34, hospital.Mild)
	p2 := hospital.NewPatient("P-002", "Yeimy Padilla", 28, hospital.Severe)
	p3 := hospital.NewPatient("P-003", "Rafael Ortega", 52, hospital.Moderate)
	p4 := hospital.NewPatient("P-004", "Wilfrido Berrio", 61, hospital.Severe)
	p5 := hospital.NewPatient("P-005", "Dayana Salcedo", 19, hospital.Mild)
	for _, p := range []*hospital.Patient{p1, p2, p3, p4, p5} {
		check(h.AdmitPatient(p))
	}
	fmt.Printf("%s: %d doctores, %d camillero, %d habitaciones, %d pacientes\n",
		h.Name(), len(h.Doctors()), len(h.Orderlies()), len(h.Rooms()), len(h.Patients()))

	attacks := []struct {
		patient  *hospital.Patient
		location string
	}{
		{p1, "cafeteria"},
		{p2, "pasillo 1, primer piso"},
		{p3, "radiology queue"},
		{p4, "hallway 2, second floor"},
	}
	var firstFailure error
	fmt.Println("\n== Ataques de sueño ==")
	for _, a := range attacks {
		rec, err := h.HandleSleepAttack(a.patient, a.location)
		switch {
		case err == nil:
			fmt.Printf("  %s %s -> room %d (atendió %s)\n", a.patient.ID(), a.location, h.RoomOf(a.patient).Number(), rec.AttenderID())
		case errors.Is(err, hospital.ErrNoRoomAvailable):
			firstFailure = err
			fmt.Printf("  %s %s -> sin habitación (atendió %s)\n  ERROR: %v\n", a.patient.ID(), a.location, rec.AttenderID(), err)
		default:
			log.Fatalf("error inesperado: %v", err)
		}
	}
	printHallway(h, "Patients asleep in hallway (antes de despertar a nadie)")

	check(h.WakePatient(p1))
	fmt.Printf("\n%s (%s) despertó y liberó su habitación\n", p1.ID(), p1.Name())
	room, err := h.AssignRoom(p4)
	check(err)

	fmt.Println("\n===== Consultas de la sección 5 =====")
	printHallway(h, "5.1 Patients asleep in hallway")

	for _, d := range h.Doctors() {
		fmt.Printf("\n== 5.2 Episodes attended by Dr. %s ==\n", d.Name())
		for _, e := range h.EpisodesByAttender(d.ID()) {
			fmt.Printf("  %s %s %s\n", e.Timestamp().Format("2006-01-02 15:04"), e.PatientID(), e.Location())
		}
	}

	fmt.Printf("\n== 5.3 Assigning room to %s ==\n", p4.ID())
	fmt.Printf("  primer intento: %v\n", firstFailure)
	fmt.Printf("  tras liberarse una cama: asignada a room %d\n", room.Number())

	fmt.Println("\n== 5.4 Severity report (Severe patients) ==")
	report := h.SevereReport()
	severe := make([]*hospital.Patient, 0, len(report))
	for p := range report {
		severe = append(severe, p)
	}
	sort.Slice(severe, func(i, j int) bool { return severe[i].ID() < severe[j].ID() })
	for _, p := range severe {
		fmt.Printf("  %s %-16s %d episodes today\n", p.ID(), p.Name(), report[p])
	}
}

func printHallway(h *hospital.Hospital, title string) {
	fmt.Printf("\n== %s ==\n", title)
	patients := h.PatientsInHallway()
	if len(patients) == 0 {
		fmt.Println("  (nadie)")
	}
	for _, p := range patients {
		fmt.Printf("  %s %-16s (%s)\n", p.ID(), p.Name(), p.Location())
	}
}

func check(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func must[T any](v T, err error) T {
	check(err)
	return v
}
