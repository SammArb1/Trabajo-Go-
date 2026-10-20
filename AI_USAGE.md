# AI_USAGE.md — Hospital de los Costeños con Narcolepsia

Universidad EIA · Lenguajes y Compiladores · Go 1.21+

Este documento describe, con honestidad, cómo usamos herramientas de IA en el
proyecto, qué errores encontramos en el código generado y qué aprendió cada
integrante.

## 1. Integrantes y reparto del trabajo

| Integrante | Rama | Responsabilidad |
|---|---|---|
| 1 — Samuel | `feature/patient-room` | Base del dominio: `types.go`, `person.go`, `patient.go`, `room.go` y las pruebas unitarias de habitaciones (`hospital_test.go`). |
| 2 | `feature/staff-episodes` | Personal y episodios: `attender.go`, `doctor.go`, `orderly.go`, `episode.go`, `attenders_test.go` y este documento. |
| 3 | `feature/hospital-orchestrator` | Orquestación: `hospital.go` (consultas y asignación de camas), `main.go` (escenario obligatorio) y `README.md`. |

## 2. Herramientas usadas y para qué

| Herramienta | Uso |
|---|---|
| **Claude** | Analizar el enunciado, repartir el trabajo, diseñar las estructuras, completar código, escribir pruebas y revisar el proyecto completo. |
| **Go toolchain** (`gofmt`, `go vet`, `go test`) | Verificar formato, advertencias y que las pruebas pasen. Es la validación real; no confiamos solo en lo que dice la IA. |
| **Git** | Una rama por integrante y unión final en `main`. |

## 3. Prompts decisivos

### Prompt 1 — Extraer la información clave y repartir el trabajo en partes iguales

> Se le entregó el PDF del enunciado y se le pidió extraer toda la información
> valiosa del proyecto (entidades, reglas técnicas, consultas, escenario y
> rúbrica) y distribuirla entre tres integrantes en partes iguales de trabajo.

**Por qué fue decisivo:** convirtió un enunciado largo en tres paquetes de
trabajo balanceados y con archivos separados, lo que evitó conflictos de
`git` y permitió trabajar en paralelo. De ahí salió el reparto de la
sección 1.

### Prompt 2 — Preguntar por las mejores opciones de diseño

> Se le pidió comparar alternativas de diseño: cómo modelar el personal para
> usar polimorfismo, dónde guardar el historial de episodios y qué tipos de
> retorno usar en las consultas.

**Por qué fue decisivo:** las decisiones que tomamos a partir de esa
conversación fueron:

- `Person` se **embebe por valor** en `Patient`, `Doctor` y `Orderly`
  (composición, sin repetir campos).
- `Doctor` y `Orderly` implementan `Attender`, y el hospital guarda su
  personal en un `[]Attender`.
- Cada doctor guarda su propio historial para que `MyEpisodes()` no necesite
  recibir el hospital.
- Los estados son tipos con constantes (`iota`), no strings sueltos.
- Los métodos que pueden fallar devuelven `error`; no se usa `panic`.

### Prompt 3 — Revisar que todo funcionara correctamente

> Se le entregó el repositorio completo y se le pidió revisar que todo
> funcionara: que compilara, que `go vet` y `go test ./...` pasaran, que se
> cumplieran los requisitos de la lista de verificación del enunciado y que
> las consultas devolvieran los resultados esperados.

**Por qué fue decisivo:** la revisión encontró diferencias entre el código y
lo que pide el PDF (ver sección 4). También detectó que `gofmt` marcaba
archivos por saltos de línea de Windows (CRLF), lo que se corrige con
`gofmt -w .`.

## 4. Casos donde el código generado estaba mal, incompleto o no era idiomático

1. **`Doctor` no guardaba sus episodios.** La primera versión solo agregaba el
   episodio al paciente, por lo que `MyEpisodes()` (consulta 5.2) no se podía
   responder sin pasarle el hospital al doctor, algo que el enunciado prohíbe.
   *Solución:* el doctor mantiene su propio slice `episodes`, y
   `MyEpisodes()` devuelve una **copia** para no exponer el estado interno.

2. **IDs de episodio con `time.Now().UnixNano()`.** Funcionaban, pero eran
   largos, dependían del reloj y dos llamadas muy seguidas podrían coincidir.
   *Solución:* un contador `atomic.Int64` que genera `EP-0001`, `EP-0002`…,
   seguro incluso si más adelante se usan goroutines.

3. **Código duplicado y sin validación en `Attend`.** `Doctor` y `Orderly`
   repetían la misma lógica y permitían "atender" a un paciente despierto.
   *Solución:* una función auxiliar compartida `newAttendRecord` y un error de
   dominio `ErrPatientAwake`.

4. **`EpisodeRecord` incompleto.** Faltaban `Summary()` y la habitación
   asignada que pide el enunciado. *Solución:* campo `room`, método `WithRoom`
   (el registro es inmutable: devuelve una copia) y `Summary()` con fecha,
   paciente, lugar y habitación.

5. **`Hospital` no seguía la especificación.** Guardaba el personal en mapas
   separados (`doctors`, `orderlies`) en lugar de un `[]Attender`, y los
   métodos se llamaban `RegisterPatient`/`RegisterDoctor` en vez de
   `AdmitPatient`/`HireDoctor`. Con mapas separados, la interfaz `Attender`
   no se usaba polimórficamente y ese ítem puntúa cero. *Solución:* guardar
   todo el personal en un `[]Attender`, renombrar los métodos según el PDF e
   implementar `PatientsInHallway`, `AssignRoom` y `SevereReport`.

6. **El paciente no sabía en qué habitación estaba.** Sin ese dato, el
   episodio no podía registrar la habitación. *Solución:* se agregó el campo
   `room` con `Room()` y `SetRoom()` en `Patient`, y `AssignRoom` lo actualiza.

## 5. Qué aprendió cada integrante

**Integrante 1 — Samuel**
- Cómo se modela encapsulamiento en Go: campos en minúscula (no exportados)
  accesibles solo mediante métodos, y por qué el paquete es la frontera.
- Cómo declarar tipos enumerados con `iota` y `String()` para que el
  compilador detecte errores que con strings pasarían desapercibidos.
- Que `Room.Occupy` y `Release` deben validar el cupo y devolver `error` en
  lugar de usar `panic`.

**Integrante 2**
- Que el polimorfismo en Go se logra con interfaces satisfechas de forma
  implícita: `Doctor` y `Orderly` cumplen `Attender` sin declararlo, y
  `var _ Attender = (*Doctor)(nil)` lo verifica en compilación.
- La diferencia entre embeber `Person` por valor y por puntero, y por qué
  elegimos valor (no hay `nil`, y la copia es independiente).
- Por qué devolver copias de slices protege el encapsulamiento, y por qué los
  métodos que modifican estado usan receptor de puntero.

**Integrante 3**
- Cómo un agregado (`Hospital`) coordina varias entidades sin conocer el tipo
  concreto del personal, usando `[]Attender`.
- Por qué `AssignRoom` devuelve `(*Room, error)` y cómo el programa maneja la
  falta de camas sin caerse, dejando al paciente en el pasillo.
- Por qué `SevereReport` devuelve `map[*Patient]int` y no imprime: la
  lógica de negocio queda en el paquete `hospital` y `main` solo imprime.

## 6. Declaración de honestidad

Usamos IA como apoyo para diseñar, escribir y revisar código, pero cada
integrante revisó, probó y entiende lo que entregó. Verificamos el resultado
con `gofmt`, `go vet` y `go test ./...`, y estamos en capacidad de explicar y
modificar cualquier parte del proyecto en la sustentación.