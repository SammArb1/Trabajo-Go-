# Hospital de los Costeños con Narcolepsia

Ejercicio de Programación Orientada a Objetos (Universidad EIA) en Go 1.21+, solo con la librería estándar.

**Integrantes:** 
Samuel Arbelaez
Nicolas Aguas
Jeronimo Aristizabal

## Cómo ejecutar

```bash
go run .          # escenario de demostración y las 4 consultas
go test ./...     # pruebas
gofmt -l .        # no debe listar archivos
go vet ./...      # no debe mostrar nada
```

## Estructura

```
main.go          # escenario e impresión
hospital/        # modelo: Person, Patient, Doctor, Orderly, Room, EpisodeRecord y Hospital
```

## Decisiones de diseño

- **Personal polimórfico:** `Hospital` guarda a los doctores y camilleros juntos en un `[]Attender`. `HandleSleepAttack` elige al siguiente de forma rotativa y llama a `Attend` sin saber qué tipo concreto es.
- **Mapas y listados ordenados:** pacientes, habitaciones y doctores se guardan en mapas por id para buscarlos directo. Como el orden de un mapa en Go es aleatorio, los listados se devuelven ordenados; así "la primera habitación libre" es siempre la de menor número.
- **Errores en vez de panic:** los casos previsibles devuelven `error`. "No hay cama" es el error `ErrNoRoomAvailable` y se detecta con `errors.Is`.
- **Sin cama:** `HandleSleepAttack` registra el episodio igual y devuelve el error; el paciente sigue en `AsleepInHallway`.
- **`SevereReport` devuelve `map[*Patient]int`:** es una relación paciente → número de episodios del día.
- **Consultas que devuelven datos:** `PatientsInHallway`, `AssignRoom`, `EpisodesByAttender` y `SevereReport` no imprimen; la impresión está solo en `main`.
- **Paquetes:** `hospital` tiene la lógica y `main` solo arma el escenario e imprime.
