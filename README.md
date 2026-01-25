# Watcher Service

Serviço de monitoramento de sistema de arquivos em tempo real, parte do projeto **What Do I Study?** — um tracker automático de estudo de programação.

## Visão Geral

O Watcher é o **sensor** do sistema. Ele monitora diretórios de estudo e detecta alterações em arquivos, gerando eventos estruturados que são publicados via HTTP e mensageria (RabbitMQ).

```
[ Watcher ] → [ Collector ] → [ Analyzer ]
     ↑
  você está aqui
```

### Papel no Sistema

- **Observa** alterações no filesystem em tempo real
- **Publica** eventos estruturados
- **Não decide nada** — apenas coleta e repassa

## Funcionalidades

- Monitoramento recursivo de diretórios
- Detecção de eventos: `CREATE`, `MODIFY`, `REMOVE`, `RENAME`, `CHMOD`
- Filtragem automática de diretórios de ruído
- **Debounce** (750ms) para evitar eventos duplicados
- Publicação via HTTP (Fiber) e RabbitMQ
- Armazenamento in-memory para consulta via API

## Exemplo de Evento

```json
{
  "timestamp": "2025-09-06T20:44:53Z",
  "file_path": "/app/dev/studies/node-rocketseat/server.js",
  "extension": ".js",
  "event": "MODIFY"
}
```

## Arquitetura

O projeto segue **Arquitetura Hexagonal** (Clean Architecture) com separação clara de responsabilidades:

```
cmd/
└── api/
    └── main.go                  # Entry point

internal/
├── domain/
│   ├── fileevent.go             # Entidade de domínio + interface Publisher
│   ├── repository/
│   │   └── memory/
│   │       └── memory.go        # Repositório in-memory
│   └── service/
│       └── watcher/
│           └── watcher.go       # Lógica de negócio (debounce, publicação)
│
└── adapter/
    ├── fsnotify/
    │   ├── fsnotify.go          # Monitoramento do filesystem
    │   └── utils.go             # Filtros de diretórios
    ├── consolepub/
    │   └── consolepub.go        # Publisher para console (debug)
    ├── rabbitmq/
    │   ├── publisher.go         # Publisher para RabbitMQ
    │   ├── connection.go        # Gerenciamento de conexão AMQP
    │   └── config.go            # Configuração via env vars
    └── httppub/
        ├── httppub.go           # API REST (Fiber)
        └── dto.go               # DTOs de resposta
```

### Fluxo de Eventos

```
┌─────────────────┐
│  fsnotify       │  detecta mudanças no filesystem
└────────┬────────┘
         ↓
┌─────────────────┐
│ WatcherService  │  aplica debounce → salva no repositório → publica
└────────┬────────┘
         ↓
    ┌────┼────┐
    ↓    ↓    ↓
┌───────┐ ┌────────┐ ┌──────┐
│Console│ │RabbitMQ│ │ HTTP │
└───────┘ └────────┘ └──────┘
```

## API HTTP

| Endpoint       | Método | Descrição                    |
|----------------|--------|------------------------------|
| `/ping`        | GET    | Health check                 |
| `/events`      | GET    | Lista todos os eventos       |
| `/events/:id`  | GET    | Busca evento por ID          |

**Porta:** `3000`

## Configuração

### Variáveis de Ambiente

| Variável            | Descrição                    | Default                              |
|---------------------|------------------------------|--------------------------------------|
| `WATCH_PATH`        | Diretório a monitorar        | `.` (diretório atual)                |
| `RABBITMQ_URI`      | URI de conexão RabbitMQ      | `amqp://guest:guest@localhost:5672/` |
| `RABBITMQ_EXCHANGE` | Nome do exchange             | `file_events`                        |
| `RABBITMQ_QUEUE`    | Nome da fila                 | `file_events_queue`                  |

### Diretórios Ignorados

O watcher ignora automaticamente:
- `node_modules`
- `.git`
- `vendor`
- `dist`

## Instalação e Execução

### Pré-requisitos

- Go 1.24.6+
- Docker (opcional)
- RabbitMQ (opcional, para mensageria)

### Desenvolvimento Local

```bash
# Clonar
git clone https://github.com/leonardoTavaresM/watcher.git
cd watcher

# Instalar dependências
go mod download

# Rodar
go run cmd/api/main.go
```

### Com Docker

```bash
# Build e run
make build
make run

# Ou com docker-compose (inclui RabbitMQ)
docker-compose up
```

### Makefile

```bash
make build    # Compila a aplicação
make run      # Executa o container
make clean    # Remove o container
make rebuild  # Limpa, builda e roda
```

## Dependências

| Biblioteca                                                        | Uso                          |
|-------------------------------------------------------------------|------------------------------|
| [fsnotify](https://github.com/fsnotify/fsnotify)                  | Monitoramento de filesystem  |
| [Fiber](https://github.com/gofiber/fiber)                         | Web framework                |
| [amqp091-go](https://github.com/rabbitmq/amqp091-go)              | Cliente RabbitMQ             |
| [uuid](https://github.com/google/uuid)                            | Geração de UUIDs             |

## Decisões Técnicas

### Por que Golang?

- Baixo consumo de recursos
- Concorrência nativa (goroutines)
- Ideal para serviços de monitoramento

### Por que Arquitetura Hexagonal?

- Domínio isolado de detalhes de implementação
- Adapters pluggáveis (fácil adicionar novos publishers)
- Testabilidade

### Por que Mensageria?

- Desacopla os serviços
- Collector pode consumir eventos de forma assíncrona
- Resiliência em caso de falhas

## Adicionando Novos Publishers

Implemente a interface `Publisher`:

```go
type Publisher interface {
    Publish(event FileEvent) error
    Close() error
}
```

Exemplo:

```go
type WebhookPublisher struct {
    url string
}

func (p *WebhookPublisher) Publish(event domain.FileEvent) error {
    // Envia evento via HTTP POST
    return nil
}

func (p *WebhookPublisher) Close() error {
    return nil
}
```

## Próximos Passos do Projeto

Este serviço é o primeiro componente do **What Do I Study?**. Os próximos serviços são:

| Serviço       | Stack           | Responsabilidade                                      |
|---------------|-----------------|-------------------------------------------------------|
| **Collector** | Go + SQL        | Consome eventos, valida, persiste em banco            |
| **Analyzer**  | Python + IA     | Analisa padrões, gera insights sobre estudos          |

### O que o sistema completo responde

- "O que eu estudei essa semana?"
- "Tenho focado mais em backend ou frontend?"
- "Meu estudo é mais noturno ou diurno?"
- "Estou alternando demais de contexto?"

## Licença

MIT

## Autor

Leonardo Tavares Malt - [@leonardoTavaresM](https://github.com/leonardoTavaresM)
