# PLAN DE TRANSFORMACIÓN GO-REDUX

## Respuesta breve
Los gophers **pueden aceptar** una arquitectura más estructurada **si aporta valor real**.  
Pero en Go, **simplicidad primero** y estructura adicional solo cuando el tamaño/complexidad lo justifique.

---

## Observaciones de ecosistema (resumen)
- En discusiones de la comunidad Go sobre Clean Architecture, el patrón general es:
  - “depende del tamaño del proyecto”
  - “evita capas académicas”
  - “usa `internal/` para encapsular, no para imitar Java”
- Layouts populares recomiendan:
  - `cmd/` para entry points
  - `internal/` para privacidad
  - sin prescribir obligatoriamente *domain/application/ports*.
- Mensaje común: pragmatismo > dogma.

**Conclusión operativa:**  
Tu intuición es correcta: en Go, forzar un layout `domain/application` sin necesidad puede verse como anti-idiomático.  
Para ejemplos de go-redux, mejor una estructura **feature-first**, cohesionada y plana.

---

## Visión y misión (versión corta)
**Visión:**  
Posicionar **go-redux** como una librería ligera y modular de coordinación de estado y flujos para backends distribuidos en Go.

**Misión:**  
Demostrar, con ejemplos ejecutables, que semánticas tipo Redux pueden resolver:
- sagas y compensaciones
- event sourcing + CQRS
- workflows de larga duración
- estado en tiempo real (gaming)

Todo ello sin imponer capas “académicas” y manteniendo un diseño Go-friendly.

---

## Cambios ya realizados (resumen técnico)
1. Reescritura completa a v2:
   - nuevo módulo `github.com/janmbaco/go-redux/v2`
   - dependencia única de `go-infrastructure/v2`
   - eliminación del antiguo `src/`
   - nuevos paquetes: `actions`, `handlers`, `events`, `ioc`
   - `store.go` genérico.
2. Store modular con handlers pluggables:
   - `Dispatch` recorre módulos vía reflexión
   - soporta `map[string]any` por selector
   - subscriptores como callbacks vía `EventManager`.
3. API basada en `ActionHandlerBuilder`:
   - mapea acciones a reducers
   - composición por módulos/slices de estado.
4. Ejemplos orientados a sistemas distribuidos:
   - saga, event sourcing, workflow, juego
   - wallet event-sourced completo
   - integración DI vía go-infrastructure.

---

## Enfoque que sigue
- Librería de coordinación ligera para backends distribuidos
  (sagas, event sourcing, CQRS, workflows).
- Integración IoC/DI y almacenamiento/eventos externos intercambiables.
- Promesa de arquitectura modular y DDD-friendly **en sentido pragmático**:
  cohesión por feature, reducers declarativos y composición por módulos.

---

## Riesgos y brechas detectadas
1. **Inmutabilidad comprometida**
   - `Dispatch` muta directamente `map[string]any`
   - no copia estado antes de publicarlo
   - handlers incompatibles pueden ser ignorados silenciosamente.
2. **Type safety parcial**
   - reflexión en `ActionHandlerBuilder.On`
   - firmas erróneas pueden panicar en runtime.
3. **Ejemplo wallet incompleto respecto a IoC/SQLite**
   - acoplo a `MemoryEventStore`
   - no se muestra el “swap” real de backends.
4. **Bug menor de logging en replay**
   - formateo incorrecto del número de eventos.
5. **Ventana de carrera en proyecciones**
   - desbloqueo/rebloqueo dentro de `Rebuild`.
6. **Desalineación docu ↔ código**
   - README/CHANGELOG prometen features aún no implementadas
     (middleware, selectores memoizados, etc.).

---

## Criterios globales de éxito
- README y CHANGELOG **sin promesas falsas**.
- Store con modos de ejecución claros:
  - `strict types`
  - `immutable/copy mode` (o alternativa documentada).
- Ejemplos que demuestran:
  - coordinación real
  - persistencia real intercambiable
  - ejecución local simple con Docker.
- 1–2 entradas claras para colaboradores (`good first issue`).

---

# FASE 1: LIMPIEZA RADICAL (1 hora)

**Objetivo:** eliminar todo rastro de framing CRUD/UI o DDD académico innecesario.

**Acciones:**
1. Borrar `examples/banking-ddd/`.
2. Transformar `README.md`:
   - eliminar: “Redux Pattern for DDD in Go”
   - nuevo tagline:  
     **“Redux Pattern for Distributed Systems Coordination in Go”**
   - reducir referencias a “domain/application layers”.
3. Actualizar badges y descripción del repositorio.
4. Limpiar `go.mod` y docs relacionadas.

**Criterio de éxito:**
- Cero menciones a banking/account/transaction en el repo.
- README posiciona go-redux como librería de coordinación distribuida.

---

# FASE 2: SAGA ORDER FULFILLMENT (10 horas)

**Estructura pragmática Go (feature-first):**

```

examples/saga-order-fulfillment/
├── cmd/server/main.go
├── internal/
│   ├── order/
│   │   ├── order.go
│   │   ├── saga.go
│   │   └── store.go
│   ├── payment/
│   │   ├── service.go
│   │   └── handlers.go
│   ├── inventory/
│   │   ├── service.go
│   │   └── handlers.go
│   ├── shipping/
│   │   ├── service.go
│   │   └── handlers.go
│   └── compensations/
│       └── handlers.go
├── api/http.go
├── docker-compose.yml
├── go.mod
└── README.md

```

**Características clave:**
1. Saga con 3 pasos: Payment → Inventory → Shipping.
2. Compensaciones automáticas.
3. Servicios mock con fallos aleatorios.
4. Ejecución concurrente de múltiples orders.
5. Correlation IDs en logs.
6. Demo simple con `curl`.

---

# FASE 3: EVENT-SOURCED WALLET (12 horas)

**Estructura pragmática Go:**

```

examples/event-sourced-wallet/
├── cmd/server/main.go
├── internal/
│   ├── wallet/
│   │   ├── wallet.go
│   │   ├── events.go
│   │   └── store.go
│   ├── eventstore/
│   │   ├── memory.go
│   │   └── sqlite.go
│   ├── projections/
│   │   └── balance.go
│   └── snapshots/
│       └── snapshotter.go
├── api/http.go
├── go.mod
└── README.md

```

**Características clave:**
1. Comandos: Deposit, Withdraw, Transfer.
2. Replay completo.
3. Snapshots periódicos.
4. CQRS con proyecciones.
5. Swap real de backend: Memory ↔ SQLite vía IoC.

---

# FASE 4: WORKFLOW EMPLOYEE ONBOARDING (10 horas)

**Estructura pragmática Go:**

```

examples/workflow-employee-onboarding/
├── cmd/server/main.go
├── internal/
│   ├── employee/
│   │   ├── employee.go
│   │   └── store.go
│   ├── workflow/
│   │   ├── steps.go
│   │   ├── engine.go
│   │   └── handlers.go
│   ├── tasks/
│   │   ├── hr_setup.go
│   │   ├── it_provisioning.go
│   │   ├── training.go
│   │   └── badge.go
│   └── persistence/
│       └── repository.go
├── api/http.go
├── go.mod
└── README.md

```

**Características clave:**
1. Workflow multi-step.
2. Dependencias entre pasos.
3. Rollback.
4. Persistencia y reanudación tras crash.
5. Demo con mocks y delays.

---

# FASE 5: GAME TURN-BASED STRATEGY (14 horas)

**Justificación:**
- Caso de uso natural para state management y replays.
- Valida time travel + event sourcing de forma muy visual.

**Estructura pragmática Go:**

```

examples/game-server/
├── cmd/server/main.go
├── internal/
│   ├── game/
│   │   ├── game.go
│   │   ├── board.go
│   │   └── store.go
│   ├── player/
│   │   ├── player.go
│   │   └── actions.go
│   ├── matchmaking/
│   │   ├── queue.go
│   │   └── matcher.go
│   ├── replay/
│   │   ├── recorder.go
│   │   └── player.go
│   └── anticheat/
│       └── validator.go
├── api/
│   ├── http.go
│   └── websocket.go
├── go.mod
└── README.md

```

**Características clave:**
1. Juego turn-based simplificado.
2. Matchmaking.
3. Redux state: board, turn, history.
4. Time travel y replay.
5. Anti-cheat básico.
6. Updates en tiempo real.

---

# FASE 6: DOCUMENTACIÓN GLOBAL (8 horas)

**Objetivo:** reposicionar go-redux sin prometer humo.

**README principal:**
1. Tagline nuevo.
2. Use cases reales (saga, event sourcing, workflows, gaming).
3. “When NOT to use”.
4. Filosofía de arquitectura Go:
   - estructura por feature
   - `internal/` para encapsular
   - sin capas académicas por defecto.
5. Roadmap breve.

**READMEs de ejemplos:**
- problema → enfoque → cómo ejecutar → trade-offs.

---

# FASE 7: TESTING Y DOCKER COMPOSE (8 horas)

**Testing strategy:**
1. Unit tests:
   - reducers/handlers
   - compensaciones
   - replay/snapshots
   - rules del juego.
2. Integration tests:
   - saga end-to-end
   - wallet con Memory/SQLite
   - workflow recovery.
3. Concurrency tests del Store.
4. Benchmarks razonables y honestos.

**Docker Compose:**
- 1 compose por ejemplo, con mocks y storage local.

---

# FASE 8: RELEASE Y COMMUNITY (4 horas)

**Release:**
- versión centrada en:
  - coherencia doc↔código
  - ejemplos ejecutables
  - mejoras de seguridad/concurrencia
  - guía de migración si aplica.

**Community (enfoque feedback-first):**
- Post en Reddit pidiendo opinión de API + ejemplos.
- Artículo:  
  **“Why Redux isn’t just for UIs: coordination patterns in Go backends”**

**Métricas realistas:**
- crecimiento visible de estrellas
- 1–3 PRs pequeños
- feedback útil que puedas convertir en issues.

---

## Resumen ejecutivo

**Tiempo total estimado:** ~67 horas

**Estrategia principal:**
- Reposicionar sin dogmas.
- Mostrar valor real con 4 ejemplos “productivos”.
- Asegurar credibilidad técnica:
  - inmutabilidad/locking
  - fail-fast en handlers
  - tests de concurrencia
  - doc honestas.

**Decisión de arquitectura para ejemplos:**
- Sin carpetas `domain/application` por defecto.
- Estructura por feature, cohesionada.
- `internal/` para encapsulación real.


