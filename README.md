# ExpenseTracker

ExpenseTracker é um aplicativo de rastreamento de despesas construído em Golang utilizando Clean Architecture. Ele permite que os usuários adicionem, visualizem, atualizem e excluam suas despesas, categorizando-as e visualizando relatórios detalhados.

## Visão Geral

O objetivo do ExpenseTracker é ajudar os usuários a gerenciar suas finanças pessoais de forma eficiente. Ele fornece uma API para interagir com o sistema, permitindo a criação e consulta de despesas, bem como a geração de relatórios.

## Tecnologias Utilizadas

- Go
- Gin (framework web)
- GORM (ORM)
- PostgreSQL (banco de dados)
- JWT (autenticação)
- Docker (para contêineres)

## Arquitetura

O projeto segue os princípios de Clean Architecture, dividindo o código em camadas bem definidas:

- **Domain**: Entidades de domínio e interfaces (ports)
- **Application**: Casos de uso (interactors)
- **Interface (API)**: Controladores e manipuladores de rotas
- **Infrastructure**: Repositórios e adaptadores de saída

## Configuração do Projeto

### Pré-requisitos

- [Golang](https://golang.org/doc/install) (versão 1.17 ou superior)
- [Docker](https://www.docker.com/products/docker-desktop) (para desenvolvimento e testes)
- [PostgreSQL](https://www.postgresql.org/download/) (local ou via Docker)

### Passo a Passo

1. Clone o repositório:

   ```bash
   git clone https://github.com/seu-usuario/expense-tracker.git
   cd expense-tracker
   ```

2. Configure o banco de dados PostgreSQL. Você pode usar Docker:

   ```bash
   docker run --name expense-db -e POSTGRES_PASSWORD=senha -e POSTGRES_USER=usuario -e POSTGRES_DB=expensedb -p 5432:5432 -d postgres
   ```

3. Crie um arquivo `.env` com as configurações do banco de dados:

   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=usuario
   DB_PASSWORD=senha
   DB_NAME=expensedb
   JWT_SECRET=sua_secret_key
   ```

4. Instale as dependências:

   ```bash
   go mod tidy
   ```

5. Execute as migrações do banco de dados:

   ```bash
   go run main.go migrate
   ```

6. Inicie a aplicação:
   ```bash
   go run main.go
   ```

A API estará disponível em `http://localhost:8080`.

## Endpoints da API

### Autenticação

- `POST /register`: Registra um novo usuário
- `POST /login`: Autentica um usuário e retorna um token JWT

### Despesas

- `POST /expenses`: Cria uma nova despesa
- `GET /users/:user_id/expenses`: Retorna as despesas de um usuário
- `PUT /expenses/:expense_id`: Atualiza uma despesa
- `DELETE /expenses/:expense_id`: Deleta uma despesa

## Tarefas

- [X] Configurar o ambiente de desenvolvimento
- [X] Implementar a camada de domínio
  - [X] Criar entidades (Expense, User, Category)
  - [X] Definir interfaces (ports)
- [ ] Implementar a camada de aplicação
  - [ ] Criar casos de uso (AddExpense, GetUserExpenses, UpdateUserExpense, RemoveExpense)
- [ ] Implementar a camada de interface (API)
  - [ ] Configurar rotas
  - [ ] Criar handlers
- [ ] Implementar a camada de infraestrutura
  - [ ] Configurar banco de dados PostgreSQL
  - [ ] Implementar repositórios (ExpenseRepository)
- [ ] Implementar autenticação JWT
- [ ] Escrever testes unitários
- [ ] Criar documentação da API com Swagger
- [ ] Configurar Docker para ambiente de desenvolvimento e produção
- [ ] Criar scripts de migração de banco de dados

## Contribuindo

1. Faça um fork do projeto
2. Crie uma nova branch (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas alterações (`git commit -am 'Adiciona nova funcionalidade'`)
4. Faça o push para a branch (`git push origin feature/nova-funcionalidade`)
5. Crie um novo Pull Request

## Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para mais detalhes.
