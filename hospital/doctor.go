package hospital

import "fmt"

// Doctor representa a un médico especialista del hospital.
// Embebe a Person (composición) e implementa la interfaz Attender.
type Doctor struct {
	Person
	specialty string
	patients  []*Patient      // pacientes a su cargo
	episodes  []EpisodeRecord // episodios en los que intervino
}

// Verificación en tiempo de compilación de que *Doctor satisface Attender.
var _ Attender = (*Doctor)(nil)

// NewDoctor crea un nuevo médico validando los datos requeridos.
func NewDoctor(id, name string, age int, specialty string) (*Doctor, error) {
	if id == "" {
		return nil, ErrEmptyID
	}
	if name == "" {
		return nil, ErrEmptyName
	}
	if age < 0 {
		return nil, ErrInvalidAge
	}
	if specialty == "" {
		specialty = "Medicina General / Somnología"
	}
	return &Doctor{
		Person:    NewPerson(id, name, age),
		specialty: specialty,
	}, nil
}

// Specialty retorna la especialidad médica del doctor.
func (d *Doctor) Specialty() string {
	if d == nil {
		return ""
	}
	return d.specialty
}

// SetSpecialty actualiza la especialidad médica previa validación.
func (d *Doctor) SetSpecialty(specialty string) error {
	if d == nil {
		return ErrAttenderNotFound
	}
	if specialty == "" {
		return ErrEmptyName
	}
	d.specialty = specialty
	return nil
}

// DiagnosePatient toma al paciente a cargo del doctor: lo agrega a su lista
// (sin duplicados) y lo registra como su médico asignado.
func (d *Doctor) DiagnosePatient(p *Patient) error {
	if d == nil {
		return ErrAttenderNotFound
	}
	if p == nil {
		return ErrNilPatient
	}
	for _, existing := range d.patients {
		if existing.ID() == p.ID() {
			p.SetAssignedDoctor(d)
			return nil
		}
	}
	d.patients = append(d.patients, p)
	p.SetAssignedDoctor(d)
	return nil
}

// PatientsUnderCare retorna una copia de los pacientes a cargo del doctor.
func (d *Doctor) PatientsUnderCare() []*Patient {
	if d == nil {
		return nil
	}
	out := make([]*Patient, len(d.patients))
	copy(out, d.patients)
	return out
}

// Attend atiende a un paciente dormido, registra el episodio en el historial
// del doctor y del paciente. Cumple la interfaz Attender.
func (d *Doctor) Attend(p *Patient, location string) (EpisodeRecord, error) {
	if d == nil {
		return EpisodeRecord{}, ErrAttenderNotFound
	}
	notes := fmt.Sprintf("Atención médica por Dr. %s (%s).", d.Name(), d.specialty)
	rec, err := newAttendRecord(d, p, location, notes)
	if err != nil {
		return EpisodeRecord{}, err
	}
	d.episodes = append(d.episodes, rec)
	return rec, nil
}

// AttendSleepEmergency atiende la emergencia en el lugar donde el paciente cayó dormido.
func (d *Doctor) AttendSleepEmergency(p *Patient) error {
	if p == nil {
		return ErrNilPatient
	}
	_, err := d.Attend(p, p.Location())
	return err
}

// MyEpisodes retorna una copia de todos los episodios en que intervino el doctor.
// No recibe el hospital: la relación está modelada dentro de la entidad.
func (d *Doctor) MyEpisodes() []EpisodeRecord {
	if d == nil {
		return nil
	}
	out := make([]EpisodeRecord, len(d.episodes))
	copy(out, d.episodes)
	return out
}