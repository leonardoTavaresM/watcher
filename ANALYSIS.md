# ANALYSIS.md — Contexto para Agentes Claude

> Este arquivo é uma análise gerada por um agente Claude para dar contexto rápido a
> futuros agentes que trabalharem neste repositório. Complementa o `README.md`
> (que é a documentação "oficial" e voltada ao usuário/leitor humano). Se algo
> aqui divergir do código, o código manda — atualize este arquivo.

## O que é este projeto

`watcher` é o **primeiro componente** de um sistema maior chamado **"What Do I
Study?"**, um projeto pessoal de estudo (não é produção, não é para uso real —
é um projeto de aprendizado de Go, arquitetura hexagonal e sistemas distribuídos
simples). O dono do repo é Leonardo Tavares Malt (github: leonardoTavaresM).

Ideia do sistema completo: monitorar automaticamente os diretórios de estudo do
usuário (`/home/leonardomalt/Documents/dev`) e, a partir dos eventos de
filesystem, inferir padrões de estudo — respondendo perguntas como "o que
estudei essa semana", "estou mais em backend ou frontend", "estudo mais de dia
ou de noite", "estou trocando de contexto demais".

Pipeline planejado (só o primeiro está implementado):

```
[ Watcher ] → [ Collector ] → [ Analyzer ]
   (Go)         (Go+SQL)      (Python+IA)
```

- **Watcher** (este repo): sensor. Observa o filesystem via `fsnotify`, aplica
  debounce, publica eventos (console/RabbitMQ/HTTP). Não decide nada, não tem IA.
- **Collector** (não implementado ainda): consumiria os eventos do RabbitMQ,
  validaria e persistiria em banco SQL.
- **Analyzer** (não implementado ainda): analisaria os dados persistidos com IA
  para gerar insights.

Se o usuário pedir para trabalhar em "Collector" ou "Analyzer", esses projetos
provavelmente ainda não existem neste diretório (`/what-do-i-study`) — vale
perguntar/checar antes de assumir que existem.

## Estado do repositório

- `watcher/` é **seu próprio repositório git** (tem `.git` próprio), mesmo que
  o diretório pai `what-do-i-study/` não seja um repo git. Remote:
  `https://github.com/leonardoTavaresM/watcher.git`.
- 11 commits no histórico, todos incrementais evoluindo de um protótipo simples
  para a arquitetura hexagonal atual (o commit mais recente, `c5ac99a`, foi
  justamente a migração para essa estrutura).
- Não há `.gitignore`. O binário compilado `main` (10MB, ELF executável) está
  no diretório raiz e **não está trackeado** pelo git (aparece como untracked
  em `git status`), mas deveria ganhar um `.gitignore` para evitar que alguém
  o commite por acidente.
- `teste.md` é um arquivo vazio (0 bytes) na raiz — parece resíduo/rascunho
  sem uso, não referenciado em nenhum lugar do código ou README.
- Não há testes automatizados (`*_test.go`) no projeto ainda.

## Stack técnica

- **Go 1.24.6**
- [`fsnotify`](https://github.com/fsnotify/fsnotify) — monitoramento de filesystem (inotify no Linux)
- [`gofiber/fiber v2`](https://github.com/gofiber/fiber) — web framework para a API HTTP
- [`rabbitmq/amqp091-go`](https://github.com/rabbitmq/amqp091-go) — cliente RabbitMQ
- [`google/uuid`](https://github.com/google/uuid) — dependência indireta (declarada no go.mod mas não vi uso direto no código lido — vale confirmar com `grep -r uuid internal/` antes de assumir onde é usado)
- Docker + docker-compose (sobe RabbitMQ com management UI na 15672 + o próprio watcher)

## Arquitetura (Hexagonal / Ports & Adapters)

A estrutura real de pastas (note que diverge um pouco da que está descrita no
`README.md` — o README ainda reflete uma versão anterior à última refatoração):

```
cmd/api/main.go                                  # entry point / composition root

internal/
├── domain/
│   └── entity/
│       └── fileevent.go                         # struct FileEvent (sem lógica)
│
├── application/
│   ├── port/
│   │   ├── publisher.go                         # interface Publisher
│   │   └── repository.go                        # interface EventRepository
│   └── service/
│       └── watcher_service.go                   # WatcherService: debounce + orquestração
│
└── infrastructure/
    ├── adapter/
    │   ├── fsnotify/
    │   │   ├── adapter.go                        # loop de eventos fsnotify, watch recursivo
    │   │   └── utils.go                           # ShouldIgnore (filtro de diretórios)
    │   ├── publisher/
    │   │   ├── console.go                         # publisher de debug (stdout)
    │   │   ├── rabbitmq.go                        # publisher RabbitMQ (só declara o exchange, producer puro)
    │   │   └── rabbitmq_connection.go              # conexão/canal AMQP
    │   └── repository/
    │       └── memory.go                           # repositório in-memory (map[int]FileEvent)
    ├── config/
    │   └── rabbitmq.go                             # leitura de env vars com defaults
    └── handler/
        ├── dto.go                                   # DTOs de resposta HTTP
        └── http.go                                  # handlers Fiber (GET /events, GET /events/:id)
```

Fluxo de uma mudança de arquivo:

```
fsnotify (adapter) → WatcherService.HandleFileEvent
                        → ShouldProcess (debounce 750ms por path, com mutex)
                        → repository.Save (in-memory)
                        → itera sobre publishers ([]port.Publisher) e chama Publish em cada um
                             ├── ConsolePublisher  (stdout, debug)
                             └── RabbitMQPublisher (exchange "file_events", topic)
```

A API HTTP (`GET /ping`, `GET /events`, `GET /events/:id`) roda em paralelo
(goroutine própria) e apenas **lê** do mesmo repositório in-memory que o
watcher escreve — ver seção "Pontos de atenção" abaixo sobre concorrência.

`main.go` é o composition root: instancia repo → publishers → service →
adapter fsnotify → handler HTTP → registra rotas → sobe duas goroutines
(HTTP server e fsnotify loop) → espera SIGINT/SIGTERM → `app.Shutdown()`.

## Configuração (env vars, ver `.env` / `docker-compose.yml`)

| Variável            | Default                                | Uso |
|---------------------|-----------------------------------------|-----|
| `WATCH_PATH`        | `.` (fallback se vazio)                 | Diretório raiz monitorado (recursivo) |
| `RABBITMQ_URI`       | `amqp://guest:guest@localhost:5672/`   | Conexão AMQP |
| `RABBITMQ_EXCHANGE`  | `file_events`                          | Exchange topic no RabbitMQ |
| `RABBITMQ_QUEUE`     | `file_events_queue`                    | Fila ligada ao exchange |

No `docker-compose.yml`, o container do watcher monta
`/home/leonardomalt/Documents/dev` (host, hardcoded, específico desta máquina)
em `/app/dev`, com `WATCH_PATH=/app/dev` — ou seja, em Docker ele observa
**todos** os projetos de dev do usuário, não só este.

Diretórios sempre ignorados (hardcoded em `fsnotify/utils.go`):
`node_modules`, `.git`, `vendor`, `dist`.

## Pontos de atenção / bugs conhecidos (úteis para não "descobrir de novo")

Achados ao ler o código-fonte diretamente (não estão documentados no README):

1. **Race condition no `InMemoryRepository`** (`infrastructure/adapter/repository/memory.go`):
   o map `events` é acessado sem lock. A goroutine do fsnotify escreve
   (`Save`) enquanto a goroutine HTTP lê (`GetAll`/`GetByID`) concorrentemente
   — isso é uma condição de corrida real em Go (maps não são safe para
   concorrência). Se o objetivo de estudo for "aprender sobre concorrência em
   Go", este é um excelente ponto de partida para adicionar um `sync.RWMutex`
   e discutir o porquê.

2. **IDs instáveis em `InMemoryRepository.Save`**: a chave usada é
   `len(m.events)`. Como não há `Delete` sendo chamado em lugar nenhum do
   fluxo normal, isso não quebra hoje — mas se `Delete` for usado, o próximo
   `Save` pode sobrescrever um ID existente (porque `len` diminui). Não é um
   bug ativo, mas é uma armadilha se alguém for expor um endpoint de delete.

3. **`ConsolePublisher.Publish`** ignora o parâmetro `event` recebido e sempre
   imprime `repository.GetAll()` inteiro a cada evento — funciona para debug,
   mas é O(n) a cada evento e crescerá conforme o número de eventos aumenta.

4. **`GET /ping`** retorna `c.JSON(\`{pong}\`)`, que serializa a *string*
   literal `"{pong}"` como JSON (`"\"{pong}\""`), não um objeto JSON
   `{"pong": true}` como o nome sugere. Cosmético, mas pode confundir quem
   testar o endpoint esperando um objeto.

5. **`RabbitMQPublisher`**: usa exchange do tipo `"topic"` mas a routing key
   (`RABBITMQ_ROUTING_KEY`, default `file_events_queue`) é fixa, não um
   padrão de roteamento por tipo de evento (ex: `file.created`,
   `file.modified`). Funciona, mas não aproveita o roteamento por tópico —
   seria natural evoluir para routing keys como `file.<EVENT>` se algum
   consumidor precisar filtrar seletivamente por tipo de evento.

   **Correção já aplicada:** originalmente este publisher também declarava
   e fazia bind da própria fila (`QueueDeclare`/`QueueBind`), mesmo sendo um
   producer puro — isso criava uma fila (`file_events_queue`) que nunca
   tinha consumidor nenhum e acumulava mensagens sem limite para sempre.
   Foi removido: o watcher agora só declara o exchange; declarar/bindar
   filas é responsabilidade de quem consome (ver `collector/internal/
   infrastructure/adapter/rabbitmq/consumer.go`, que faz exatamente isso
   com a fila `collector_queue`).

6. **README desatualizado em relação à estrutura de pastas**: o README
   descreve `internal/domain/fileevent.go`, `internal/domain/repository/memory/`,
   `internal/domain/service/watcher/`, `internal/adapter/...` — mas a estrutura
   real (pós-refatoração hexagonal do último commit) é
   `internal/domain/entity/`, `internal/application/{port,service}/`,
   `internal/infrastructure/adapter/...`. Ao orientar alguém pelo README,
   prefira sempre verificar a árvore real de arquivos.

7. `teste.md` (vazio) e o binário `main` não versionado na raiz parecem
   artefatos soltos — provavelmente seguros para limpar, mas confirme com o
   usuário antes de deletar (podem ser rascunho em andamento).

## Como rodar (resumo — ver README para detalhes)

```bash
# Local, sem Docker (requer RabbitMQ rodando à parte)
go run cmd/api/main.go

# Fluxo recomendado: sobe só o RabbitMQ via compose, roda o Go local
make dev

# Tudo em Docker (watcher + RabbitMQ)
docker-compose up
```

API HTTP fica em `http://localhost:3000`. RabbitMQ management UI em
`http://localhost:15672` (guest/guest).

## Sinalizadores para próximas conversas

- Se o usuário mencionar "Collector" ou "Analyzer", são os próximos serviços
  planejados do sistema "What Do I Study?" — ainda não existem no disco (até
  onde esta análise verificou). Pergunte onde ele quer criá-los.
- Este é um projeto de **estudo pessoal**, não produção — não superengenheirar,
  não adicionar abstrações além do necessário, é razoável manter simplicidade
  didática (alinhado com CLAUDE.md global de evitar over-engineering).
- Ao propor correções dos bugs listados acima, trate-os como *oportunidades de
  aprendizado* (é esse o propósito do projeto) — vale explicar o "porquê" ao
  corrigir, não só aplicar o fix silenciosamente.
