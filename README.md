# core-banking-go

API de core bancário escrita em Go, com ledger de partidas dobradas, transferências atômicas e idempotentes, autenticação OIDC, Clean Architecture, CI/CD e deploy na AWS.

![Status](https://img.shields.io/badge/status-em%20constru%C3%A7%C3%A3o-orange)
![Go](https://img.shields.io/badge/go-1.26-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-blue)

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

Em construção. Esta seção será atualizada conforme o ambiente local for configurado.

Pré-requisitos:

- [Go 1.26+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) e Docker Compose

## Fluxo de trabalho

- **Branches:** [Git Flow](https://nvie.com/posts/a-successful-git-branching-model/). `master` é produção e `develop` é a branch de integração. O desenvolvimento acontece em `feature/*` e `bugfix/*`, as versões passam por `release/*` e as correções urgentes por `hotfix/*`.
- **Commits:** [Conventional Commits](https://www.conventionalcommits.org/) (`feat`, `fix`, `chore`, `docs`, `test`, `refactor`, `ci`).
- **Pull requests:** toda mudança passa por PR com os checks aprovados. Push direto em `master` e `develop` é bloqueado.
- **Decisões:** escolhas técnicas relevantes são registradas como ADRs.

## Licença

Distribuído sob a [licença MIT](LICENSE).

