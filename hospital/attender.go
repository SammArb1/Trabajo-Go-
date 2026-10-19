package hospital

import "errors"

// ErrPatientAwake indica que se intentó atender un episodio de un paciente despierto.
var ErrPatientAwake = errors.New("el paciente está despierto: no hay episodio que atender")

// Attender define la interfaz común para cualquier personal médico o asistencial
// capaz de atender a un paciente en un episodio narcoléptico.
// Doctor y Orderly la implementan; el hospital guarda su personal en un []Attender
// y despacha a quien esté disponible sin conocer el tipo concreto (polimorfismo).
type Attender interface {
	ID() string
	Name() string
	Attend(p *Patient, location string) (EpisodeRecord, error)
}
