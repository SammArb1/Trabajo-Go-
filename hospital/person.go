package hospital

// Person representa la entidad base con atributos y métodos compartidos
// para las personas en el hospital (encapsulamiento con campos privados).
type Person struct {
	id   string
	name string
	age  int
}

// NewPerson construye una nueva instancia de Person con validación defensiva.
func NewPerson(id string, name string, age int) (*Person, error) {
	if id == "" {
		return nil, ErrEmptyID
	}
	if name == "" {
		return nil, ErrEmptyName
	}
	if age < 0 {
		return nil, ErrInvalidAge
	}
	return &Person{
		id:   id,
		name: name,
		age:  age,
	}, nil
}

// ID retorna el identificador único de la persona.
func (p *Person) ID() string {
	if p == nil {
		return ""
	}
	return p.id
}

// Name retorna el nombre completo de la persona.
func (p *Person) Name() string {
	if p == nil {
		return ""
	}
	return p.name
}

// Age retorna la edad de la persona en años.
func (p *Person) Age() int {
	if p == nil {
		return 0
	}
	return p.age
}

// SetName actualiza el nombre de la persona validando que no sea vacío.
func (p *Person) SetName(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	p.name = name
	return nil
}

// SetAge actualiza la edad de la persona validando que sea positiva.
func (p *Person) SetAge(age int) error {
	if age < 0 {
		return ErrInvalidAge
	}
	p.age = age
	return nil
}
