package hospital

// Attender define la interfaz común para cualquier personal médico o asistencial
// capaz de atender a un paciente en un episodio narcoléptico.
type Attender interface {
	ID() string
	Name() string
	Attend(p *Patient, location string) (EpisodeRecord, error)
}
