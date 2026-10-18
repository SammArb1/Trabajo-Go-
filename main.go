package main

import (
	"fmt"

	"costenos-narcolepsia/hospital"
)

func main() {
	_ = hospital.NewHospital("Hospital de los Costeños con Narcolepsia")
	fmt.Println("Hospital inicializado correctamente")
}
