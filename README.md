## Executar o Sistema

Criar um arquivo .env

Crie um arquivo .env na raiz do projeto e adicione sua chave da API do WeatherAPI:

plaintext
```sh
WEATHER_API_KEY=your_weather_api_key
```

### Executar Docker Compose

No diretório do projeto, execute:

```sh
docker-compose up --build
````

Isso irá construir as imagens Docker e iniciar os serviços definidos no docker-compose.yml. A aplicação estará disponível em http://localhost:8081 e o Zipkin em http://localhost:9411.

### Testar a Aplicação
Com a aplicação em execução, você pode testar utilizando curl ou httpie:

Teste com um CEP válido
```sh
curl -X POST -H "Content-Type: application/json" -d '{"cep": "01001000"}' http://localhost:8081/cep
```

Teste com um CEP inválido (formato incorreto)
```sh
curl -X POST -H "Content-Type: application/json" -d '{"cep": "123"}' http://localhost:8081/cep
````

Teste com um CEP não encontrado
```sh
curl -X POST -H "Content-Type: application/json" -d '{"cep": "99999999"}' http://localhost:8081/cep
````

### Visualizar Traces no Zipkin
Acesse http://localhost:9411 para visualizar os traces e spans no Zipkin.

### Executar Testes Automatizados
Para executar os testes automatizados no serviço B, utilize:

```sh
cd service_b
go test ./...
```