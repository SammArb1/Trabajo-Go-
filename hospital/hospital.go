package hospital

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

var (
	ErrNoRoomAvailable    = errors.New("no room available: patient stays in the hallway")
	ErrNoStaff            = errors.New("no hay personal contratado")
	ErrPatientNotAsleep   = errors.New("el paciente no está dormido")
	ErrPatientInBed       = errors.New("el paciente ya está en una cama")
	ErrPatientNotAdmitted = errors.New("el paciente no está admitido en el hospital")
)

// Hospital guarda el personal en un []Attender (doctores y camilleros juntos),
// y los pacientes, habitaciones y doctores en mapas por id.
type Hospital struct {
	name         string
	rooms        map[string]*Room
	patients     map[string]*Patient
	doctors      map[string]*Doctor
	orderlies    map[string]*Orderly
	staff        []Attender
	nextAttender int
	episodes     []EpisodeRecord
}

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
		staff:     make([]Attender, 0),
		episodes:  make([]EpisodeRecord, 0),
	}
}

func (h *Hospital) Name() string {
	return h.name
}

func (h *Hospital) AddRoom(room *Room) error {
	if room == nil {
		return ErrNilRoom
	}
	if _, exists := h.rooms[room.ID()]; exists {
		return errors.New("la habitación ya existe en el hospital")
	}
	h.rooms[room.ID()] = room
	return nil
}

func (h *Hospital) AdmitPatient(p *Patient) error {
	if p == nil {
		return ErrNilPatient
	}
	if _, exists := h.patients[p.ID()]; exists {
		return errors.New("el paciente ya está admitido")
	}
	h.patients[p.ID()] = p
	return nil
}

func (h *Hospital) HireDoctor(d *Doctor) error {
	if d == nil {
		return ErrAttenderNotFound
	}
	if h.isStaff(d.ID()) {
		return fmt.Errorf("el personal %s ya está contratado", d.ID())
	}
	h.doctors[d.ID()] = d
	h.staff = append(h.staff, d)
	return nil
}

func (h *Hospital) HireOrderly(o *Orderly) error {
	if o == nil {
		return ErrAttenderNotFound
	}
	if h.isStaff(o.ID()) {
		return fmt.Errorf("el personal %s ya está contratado", o.ID())
	}
	h.orderlies[o.ID()] = o
	h.staff = append(h.staff, o)
	return nil
}

func (h *Hospital) RegisterEpisode(record EpisodeRecord) error {
	if record.ID() == "" {
		return ErrEmptyID
	}
	h.episodes = append(h.episodes, record)
	return nil
}

// Los listados se devuelven ordenados porque el orden de un mapa en Go es aleatorio.

func (h *Hospital) Rooms() []*Room {
	list := make([]*Room, 0, len(h.rooms))
	for _, r := range h.rooms {
		list = append(list, r)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Number() < list[j].Number() })
	return list
}

func (h *Hospital) Patients() []*Patient {
	list := make([]*Patient, 0, len(h.patients))
	for _, p := range h.patients {
		list = append(list, p)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID() < list[j].ID() })
	return list
}

func (h *Hospital) Doctors() []*Doctor {
	list := make([]*Doctor, 0, len(h.doctors))
	for _, d := range h.doctors {
		list = append(list, d)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID() < list[j].ID() })
	return list
}

func (h *Hospital) Orderlies() []*Orderly {
	list := make([]*Orderly, 0, len(h.orderlies))
	for _, o := range h.orderlies {
		list = append(list, o)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID() < list[j].ID() })
	return list
}

func (h *Hospital) Episodes() []EpisodeRecord {
	list := make([]EpisodeRecord, len(h.episodes))
	copy(list, h.episodes)
	return list
}

// HandleSleepAttack duerme al paciente, intenta darle cama, despacha a un Attender
// y registra el episodio. Si no hay cama devuelve el episodio junto con un error que
// envuelve ErrNoRoomAvailable; el paciente queda en el pasillo.
func (h *Hospital) HandleSleepAttack(p *Patient, location string) (EpisodeRecord, error) {
	if p == nil {
		return EpisodeRecord{}, ErrNilPatient
	}
	if !h.isAdmitted(p) {
		return EpisodeRecord{}, ErrPatientNotAdmitted
	}
	if p.State() != Awake {
		return EpisodeRecord{}, ErrPatientInBed
	}
	attender, err := h.nextAttenderInTurn()
	if err != nil {
		return EpisodeRecord{}, err
	}

	p.SufferSleepAttack(location)
	_, roomErr := h.AssignRoom(p)
	if roomErr != nil && !errors.Is(roomErr, ErrNoRoomAvailable) {
		return EpisodeRecord{}, roomErr
	}

	record, err := attender.Attend(p, location)
	if err != nil {
		return EpisodeRecord{}, err
	}
	if err := h.RegisterEpisode(record); err != nil {
		return EpisodeRecord{}, err
	}
	return record, roomErr
}

// AssignRoom asigna la primera habitación libre (por número de habitación).
// Sin cama devuelve error y el paciente sigue en AsleepInHallway.
func (h *Hospital) AssignRoom(p *Patient) (*Room, error) {
	if p == nil {
		return nil, ErrNilPatient
	}
	if !h.isAdmitted(p) {
		return nil, ErrPatientNotAdmitted
	}
	switch p.State() {
	case Awake:
		return nil, ErrPatientNotAsleep
	case AsleepInBed:
		return nil, ErrPatientInBed
	}
	for _, r := range h.Rooms() {
		if r.IsAvailable() {
			if err := r.Occupy(p); err != nil {
				return nil, err
			}
			p.SetState(AsleepInBed)
			return r, nil
		}
	}
	return nil, fmt.Errorf("paciente %s: %w", p.ID(), ErrNoRoomAvailable)
}

// WakePatient despierta al paciente y libera la habitación donde estaba.
func (h *Hospital) WakePatient(p *Patient) error {
	if p == nil {
		return ErrNilPatient
	}
	room := h.RoomOf(p)
	if err := p.WakeUp(); err != nil {
		return err
	}
	if room != nil {
		return room.Release(p)
	}
	return nil
}

// RoomOf devuelve la habitación donde está el paciente, o nil si no tiene.
func (h *Hospital) RoomOf(p *Patient) *Room {
	for _, r := range h.Rooms() {
		for _, o := range r.Occupants() {
			if o == p {
				return r
			}
		}
	}
	return nil
}

// PatientsInHallway (consulta 5.1): pacientes dormidos en un pasillo.
func (h *Hospital) PatientsInHallway() []*Patient {
	var out []*Patient
	for _, p := range h.Patients() {
		if p.State() == AsleepInHallway {
			out = append(out, p)
		}
	}
	return out
}

// EpisodesByAttender (consulta 5.2): episodios en los que intervino el personal con ese id.
func (h *Hospital) EpisodesByAttender(id string) []EpisodeRecord {
	var out []EpisodeRecord
	for _, e := range h.episodes {
		if e.AttenderID() == id {
			out = append(out, e)
		}
	}
	return out
}

// SevereReport (consulta 5.4): pacientes Severe y cuántos episodios tuvieron hoy.
func (h *Hospital) SevereReport() map[*Patient]int {
	report := make(map[*Patient]int)
	severe := make(map[string]*Patient)
	for _, p := range h.patients {
		if p.Level() == Severe {
			report[p] = 0
			severe[p.ID()] = p
		}
	}
	now := time.Now()
	for _, e := range h.episodes {
		if p, ok := severe[e.PatientID()]; ok && sameDay(e.Timestamp(), now) {
			report[p]++
		}
	}
	return report
}

func sameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (h *Hospital) isAdmitted(p *Patient) bool {
	return h.patients[p.ID()] == p
}

func (h *Hospital) isStaff(id string) bool {
	for _, s := range h.staff {
		if s.ID() == id {
			return true
		}
	}
	return false
}

// nextAttenderInTurn elige al siguiente del personal de forma rotativa.
func (h *Hospital) nextAttenderInTurn() (Attender, error) {
	if len(h.staff) == 0 {
		return nil, ErrNoStaff
	}
	a := h.staff[h.nextAttender%len(h.staff)]
	h.nextAttender++
	return a, nil
}
