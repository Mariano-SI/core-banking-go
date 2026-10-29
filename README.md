# core-banking-go

API de core bancário escrita em Go, com ledger de partidas dobradas, transferências atômicas e idempotentes, autenticação OIDC, Clean Architecture, CI/CD e deploy na AWS.

![Status](https://img.shields.io/badge/status-em%20constru%C3%A7%C3%A3o-orange)
![Go](https://img.shields.io/badge/go-1.26-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-blue)
![CI](https://github.com/Mariano-SI/core-banking-go/actions/workflows/ci.yml/badge.svg?branch=develop)

## Sobre

Projeto de estudo que implementa o núcleo de um banco digital: contas, saldos, depósitos, transferências e extratos.

O domínio bancário exige que o dinheiro nunca seja criado nem perdido, que as operações permaneçam consistentes sob alta concorrência e que uma mesma requisição enviada duas vezes nunca movimente o dinheiro duas vezes. Esses requisitos tornam o projeto uma base prática para estudar arquitetura de software e system design em Go.

O projeto evolui em fases. A Fase 1 entrega uma API pronta para produção na AWS. As fases seguintes incorporam, um de cada vez, temas de system design como mensageria, cache, observabilidade, resiliência e escalabilidade.

**Este é um projeto de estudo e não deve ser usado para movimentar dinheiro real.**

## Domínio

| Conceito | Descrição |
|---|---|
| Cliente | Titular das contas, vinculado a uma identidade do provedor de identidade |
| Conta | Conta corrente cujo saldo é derivado do ledger |
| Ledger | Registro imutável de todas as movimentações, em partidas dobradas |
| Transferência | Movimentação atômica e idempotente entre duas contas |
| Depósito e saque | Movimentações com contas internas do sistema |
| Extrato | Histórico paginado dos lançamentos de uma conta |

### Regras

1. Dinheiro não é criado nem destruído, apenas movimentado. Toda transação debita uma conta e credita outra, e a soma dos lançamentos é sempre zero.
2. Lançamentos são imutáveis. Correções são feitas com novos lançamentos de estorno.
3. Toda movimentação é idempotente. Repetir uma requisição com a mesma `Idempotency-Key` não duplica a transferência.
4. Valores monetários são inteiros em centavos, nunca ponto flutuante.

### Exemplo de partidas dobradas

Uma transferência de R$ 50,00 de Alice para Bob gera uma transação com dois lançamentos:

| Transação | Conta | Débito | Crédito |
|---|---|---:|---:|
| `tx-123` | Alice | 5000 | |
| `tx-123` | Bob | | 5000 |
| | **Total** | **5000** | **5000** |

## Arquitetura

O serviço segue a Clean Architecture: as dependências apontam sempre para dentro, e as regras de negócio não conhecem HTTP, SQL ou AWS.

```
┌──────────────────────────────────────────────────────────┐
│  Adaptadores                                             │
│  Handlers HTTP, repositórios PostgreSQL, OIDC            │
│  ┌────────────────────────────────────────────────────┐  │
│  │  Casos de uso                                      │  │
│  │  AbrirConta, Depositar, Transferir, Extrato        │  │
│  │  ┌──────────────────────────────────────────────┐  │  │
│  │  │  Domínio                                     │  │  │
│  │  │  Conta, Ledger, Dinheiro, regras de negócio  │  │  │
│  │  └──────────────────────────────────────────────┘  │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

A estrutura de pastas e os limites entre as camadas serão definidos em um registro de decisão de arquitetura (ADR).

### Autenticação

A API não manipula senhas. A autenticação é delegada a um provedor de identidade OpenID Connect, e a API atua como resource server, validando apenas os tokens de acesso.

| Ambiente | Provedor de identidade |
|---|---|
| Desenvolvimento local | Keycloak (Docker) |
| Produção | AWS Cognito |

A troca de provedor é feita por configuração. A autorização, como garantir que um cliente acesse apenas as próprias contas, é uma regra de negócio da API.

## Como executar

Pré-requisitos:

- [Go 1.26+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) e Docker Compose
- `make`
- [golangci-lint](https://golangci-lint.run/) v2, para executar o linter localmente
- [goose](https://pressly.github.io/goose/), para as migrations do banco, instalado com `make tools`

### Configuração

A aplicação é configurada por variáveis de ambiente. Os valores padrão permitem executar localmente sem nenhuma configuração.

| Variável | Descrição | Padrão |
|---|---|---|
| `APP_ENV` | Ambiente: `development` ou `production` | `development` |
| `HTTP_PORT` | Porta do servidor HTTP | `8080` |
| `LOG_LEVEL` | Nível mínimo de log: `debug`, `info`, `warn` ou `error` | `info` |
| `DATABASE_URL` | String de conexão com o PostgreSQL, usada pelas migrations | — |

O arquivo `.env.example` documenta todas as variáveis. Em `production`, os logs são emitidos em JSON; em `development`, em texto.

### Ambiente local

As dependências da API rodam em containers definidos no `docker-compose.yml`:

| Serviço | Endereço | Uso |
|---|---|---|
| PostgreSQL 17 | `localhost:5432` | Banco de dados da aplicação |
| Keycloak 26 | `http://localhost:8081` | Provedor de identidade (console de administração) |

As credenciais vêm do `.env`. As senhas não têm valor padrão: sem elas, o `docker compose` não sobe.

```bash
cp .env.example .env
make up
```

Os dados do PostgreSQL ficam em um volume e são mantidos entre reinicializações. Para apagar os dados e recomeçar do zero, use `docker compose down -v`.

### Comandos

Os comandos do dia a dia estão no `Makefile`. Execute `make` para ver a lista.

| Comando | Descrição |
|---|---|
| `make run` | Executa a API, carregando as variáveis do `.env` quando o arquivo existe |
| `make build` | Gera o binário em `bin/` |
| `make test` | Executa os testes com cobertura (`make test RACE=-race` habilita o detector de race) |
| `make cover` | Executa os testes e exibe a cobertura por função |
| `make fmt` | Formata o código |
| `make vet` | Executa o `go vet` |
| `make tidy` | Sincroniza o `go.mod` e o `go.sum` com o código |
| `make clean` | Remove o binário e o arquivo de cobertura |
| `make up` | Sobe as dependências locais (PostgreSQL e Keycloak) |
| `make down` | Para as dependências locais, mantendo os dados |
| `make logs` | Acompanha os logs das dependências locais |
| `make lint` | Executa os linters (golangci-lint), conforme o `.golangci.yml` |
| `make tools` | Instala as ferramentas de desenvolvimento (goose) |
| `make migrate-up` | Aplica as migrations pendentes |
| `make migrate-down` | Desfaz a última migration aplicada |
| `make migrate-status` | Mostra quais migrations foram aplicadas |
| `make migrate-create name=<nome>` | Cria uma nova migration em `migrations/` |

Para usar variáveis próprias localmente, copie o `.env.example` para `.env` e ajuste os valores. O `.env` não é versionado.

### Migrations

O esquema do banco é versionado em `migrations/` com o [goose](https://pressly.github.io/goose/). Cada arquivo SQL tem a migration (`-- +goose Up`) e a reversão (`-- +goose Down`), numerados em sequência. O goose registra as migrations aplicadas na tabela `goose_db_version`.

```bash
make tools        # instala o goose (uma vez)
make up           # sobe o PostgreSQL
make migrate-up   # aplica as migrations pendentes
```

O esquema inicial cria as tabelas `branches`, `customers` e `accounts`, a sequence `account_number_seq` e as agências `0001` (Agência Central) e `0002` (Agência Digital).

Regras:

- Migrations já mergeadas **não são editadas**. Correções e alterações são feitas em uma migration nova.
- Toda migration tem um `Down` funcional. Antes do PR, valide o ciclo completo: `make migrate-up`, `make migrate-down` até o início e `make migrate-up` novamente.
- Em produção, as migrations rodam como etapa separada do deploy, e não na inicialização da API.

### Executando a API

```bash
make run
```

O health check fica disponível em `GET http://localhost:8080/health`.

## Fluxo de trabalho

O repositório segue o [Git Flow](https://nvie.com/posts/a-successful-git-branching-model/). As regras são aplicadas por proteção de branches e validadas automaticamente em cada pull request.

### Branches

| Branch | Função | Criada a partir de | Destino |
|---|---|---|---|
| `master` | Produção. Cada merge corresponde a uma versão | — | — |
| `develop` | Integração. Branch padrão do repositório | `master` | — |
| `feature/*` | Novas funcionalidades | `develop` | `develop` |
| `bugfix/*` | Correções durante o desenvolvimento | `develop` | `develop` |
| `release/*` | Preparação de uma versão | `develop` | `master` |
| `hotfix/*` | Correções urgentes em produção | `master` | `master` |

Merges permitidos, validados pelo workflow `check-gitflow`:

| Destino | Origens permitidas |
|---|---|
| `master` | `release/*`, `hotfix/*` |
| `develop` | `feature/*`, `bugfix/*`, `release/*`, `hotfix/*` e backport da `master` (com a label `backport`) |

### Proteção de branches

As branches `master` e `develop` são protegidas por um ruleset:

- Push direto e force push bloqueados. Toda mudança entra por pull request.
- As branches não podem ser excluídas.
- O status check `check-gitflow` é obrigatório.
- Conversas abertas no pull request precisam ser resolvidas antes do merge.

### Backport automático

Todo push na `master` (merge de `release/*` ou `hotfix/*`) dispara o workflow `backport`, que abre um pull request da `master` para a `develop` com a label `backport`. Assim, as correções feitas em produção não se perdem na branch de integração.

O pull request é criado por um GitHub App dedicado, e não pelo `GITHUB_TOKEN`. Eventos gerados pelo `GITHUB_TOKEN` não disparam outros workflows, o que deixaria o pull request de backport sem o status check obrigatório. O App tem permissões mínimas (Contents: read, Pull requests: write) e seu token é gerado a cada execução, com validade de uma hora.

### Merge

| Origem → destino | Estratégia |
|---|---|
| `feature/*`, `bugfix/*` → `develop` | Squash and merge: um commit por entrega, com o título do pull request como mensagem |
| `release/*`, `hotfix/*` → `master` | Merge commit: preserva o histórico da versão |
| Backport `master` → `develop` | Merge commit |

### Commits e pull requests

- Commits e títulos de pull request seguem o [Conventional Commits](https://www.conventionalcommits.org/) (`feat`, `fix`, `chore`, `docs`, `test`, `refactor`, `ci`).
- A descrição do pull request registra o que foi feito, as decisões tomadas e como testar.

### Decisões

Escolhas técnicas relevantes são registradas como ADRs (Architecture Decision Records).

## Licença

Distribuído sob a [licença MIT](LICENSE).
