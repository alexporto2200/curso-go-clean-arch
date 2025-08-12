# Curso Go Clean Architecture - Desafio Orders

## Tarefa
Olá devs!
Agora é a hora de botar a mão na massa. Para este desafio, você precisará criar o usecase de listagem das orders.
Esta listagem precisa ser feita com:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL
Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.
Para a criação do banco de dados, utilize o Docker (Dockerfile / docker-compose.yaml), com isso ao rodar o comando docker compose up tudo deverá subir, preparando o banco de dados.
Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.

## Portas dos Serviços
- **GraphQL Server**: 8080 (http://localhost:8080)
- **REST API**: 8081 (http://localhost:8081)
- **gRPC Server**: 8082
- **PostgreSQL**: 5432

## 🚀 Como Executar o Projeto
```bash
# Build e start completo
docker-compose up -d --build

# Ver logs
docker-compose logs -f app

# Parar serviços
docker-compose down
```

### 🧪 Testar as APIs
```bash

```

### Testar com curl
```bash
# Usar arquivo de teste
# api.http - Exemplos para REST Client, extesão do VSCode
# ou use o curl

# Listar orders rest
curl http://localhost:8081/api/v1/orders

# Criar order rest
curl -X POST http://localhost:8081/api/v1/orders -H "Content-Type: application/json" -d '{"description": "Order 1"}'

# Listar orders com gRPC
grpcurl -plaintext -proto proto/order.proto localhost:8082 order.OrderService/ListOrders

# Criar order com gRPC
grpcurl -plaintext -proto proto/order.proto -d '{"description": "Order 1"}' localhost:8082 order.OrderService/CreateOrder

# Listar orders com GraphQL
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{"query": "query { listOrders { id desc createdAt updatedAt } }"}'

# Criar order com GraphQL
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{"query": "mutation { createOrder(input: {desc: \"Nova Order via GraphQL\"}) { id desc createdAt updatedAt } }"}'


```



#### Comandos Úteis
```bash
# Apenas start (se já foi buildado)
docker-compose up -d

# Rebuild
docker-compose up -d --build

# Ver status
docker-compose ps

# Ver logs específicos
docker-compose logs -f app
docker-compose logs -f postgres

# Restart
docker-compose restart

# Limpeza completa
docker-compose down -v --remove-orphans
```




## Anotações desenvolvimento
#### 1. Configurar Banco de Dados
```bash
# Iniciar PostgreSQL
docker-compose up -d postgres

# Aguardar o banco inicializar (cerca de 10 segundos)
sleep 10
```

#### 2. Instalar Dependências
```bash
# Instalar dependências Go
go mod tidy

# Instalar protoc (se necessário)
sudo apt update && sudo apt install -y protobuf-compiler

# Instalar plugins do protoc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

#### 3. Gerar Código Protobuf
```bash
export PATH=$PATH:$(go env GOPATH)/bin
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/order.proto
```

#### 4. Executar a Aplicação
```bash
# Compilar
go build ./cmd/server

# Executar
./server
```


