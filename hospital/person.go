package hospital

// Person representa la entidad base con atributos compartidos para las personas en el hospital.
// Cumple con encapsulamiento estricto: campos no exportados y getters públicos.
type Person struct {
	id   string
	name string
	age  int
}

// NewPerson construye una nueva instancia inmutable de Person retornada por valor.
func NewPerson(id, name string, age int) Person {
	return Person{
		id:   id,
		name: name,
		age:  age,
	}
}

// ID retorna el identificador único de la persona.
func (p Person) ID() string {
	return p.id
}

// Name retorna el nombre de la persona.
func (p Person) Name() string {
	return p.name
}

// Age retorna la edad de la persona en años.
func (p Person) Age() int {
	return p.age
}
