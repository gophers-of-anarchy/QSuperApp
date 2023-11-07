# SabziPolo Bank
Sample banking system based on Sabzi Polo Mahi, written in Go.

## Features
- User registration and login
- JWT authentication
- Dockerized
- RESTful API
- CRUD operations on accounts and cards
- PostgreSQL database
    - Migrations
    - Seeders
- Redis cache
- SMS services for transfer notifications
  - Kavehnegar
  - Ghasedak
- Tests
    - Account
    - Card
    - Transfer

## Clone and run
Remember to turn on VPN if you are in Iran. Otherwise, docker will fail to download some packages.

```bash
$ git clone git@github.com:ArmanJR/sabzi-polo-bank.git
$ cd sabzi-polo-bank
$ docker-compose up
```
This will build the image, run the containers, migrate database tables, seed them with the data in `database/seeders/init_seeder.go` and start the server.
## API
### Register
```bash
curl --request POST \
  --url http://localhost:8080/register \
  --header 'Content-Type: application/json' \
  --data '{
	"username": "ali",
	"password": "daei",
	"email": "abbas@bo-azar.com",
	"cellphone": "09123456780"
}'
```

### Login
```bash
curl --request POST \
  --url http://localhost:8080/login \
  --header 'Content-Type: application/json' \
  --data '{
	"username": "ali",
	"password": "daei"
}'
```
This will return a JWT token. You need to use this token in the following requests.

### Restricted Area
```bash
curl --request GET \
  --url http://localhost:8080/restricted_area \
  --header 'Authorization: {token}'
```

### Create Account
```bash
curl --request POST \
  --url http://localhost:8080/create_account \
  --header 'Authorization: {token}' \
  --header 'Content-Type: application/json' \
  --data '{
	"name": "Savings",
	"balance": 20000
}'
```
This will return the account id which you need to use in the card creation request.

### Create Card
Feel free to use multiple card number formats (Persian, Arabic or English digits).
```bash
curl --request POST \
  --url http://localhost:8080/create_card \
  --header 'Authorization: {token}' \
  --header 'Content-Type: application/json' \
  --data '{
	"account_id": "60",
	"card_number": "۶۰۳۷-۹۹۱۱-۶۱۳۳-۳۷۸۰"
}'
```

### Transfer Money
```bash
curl --request POST \
  --url http://localhost:8080/transfer \
  --header 'Authorization: {token}' \
  --header 'Content-Type: application/json' \
  --data '{
	"from_card_number": "6037-9911-6133-3780",
	"to_card_number": "5022-2910-0000-0000",
	"amount": "15000"
}'
```
Each transaction has a `500 IRT` fee stored in `transaction_fees` table.

After a successful transfer, the server will send a notification to the both sender's & receiver's cellphone number via SMS (Update `.env` with your API keys and uncomment the code in `controllers/transaction.go:135` to enable this feature).

### Top Users
Top 3 users with the most transactions in the last 10 minutes, with their last 10 transactions.
```bash
curl --request GET \
  --url http://localhost:8080/top_users \
  --header 'Authorization: {token}'
```

## Tests
You can run the tests by running the following command:
```bash
$ docker-compose exec app go mod tidy && go test ./tests/...
```

## Design patterns used
- Singleton: in database and redis configs
- Strategy: in sms services
- Model-View-Controller
- Dependency Injection
- Middleware

## License
This project is released under the [Do What the Fuck You Want To](https://en.wikipedia.org/wiki/WTFPL) public license.