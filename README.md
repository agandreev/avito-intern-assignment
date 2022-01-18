<h1 align="center">Тестовое задание на позицию стажера-бекендера</h1>

## Libraries

- [Chi](https://github.com/go-chi/chi)
- [Swagger](https://github.com/swaggo/http-swagger)
- [Testify](https://github.com/stretchr/testify)
- [Logrus](https://github.com/sirupsen/logrus)
- [Viper](https://github.com/spf13/viper)
- [PGX](https://github.com/jackc/pgx)

## Notes

- Оба дополнительных задания выполнены.

# Problems
* В предложенном сервисе для конвертации валют в бесплатной подписке можно конвертировать валюты только в евро. Можно было бы сменить сервис на полностью бесплатный, но, чтобы не рисковать надежностью при переводе из валюты X в валюту Y я предпочел промежуточно переводить обе валюты в евро для рассчета коэффициента.
* Я не стал делать авторизацию, т.к. предположил, что в этом мире за нее отвечает другой микросервис. Однако, если это было необходимо, то в моем репозитории есть проект AlgoTrader с реализацией JWT авторизации. 

----
# Preparation

----
## Install

    git clone https://github.com/agandreev/avito-intern-assignment.git

## Fill config.env file

    API_KEY=
    DB_USER=user
    DB_PSWD=passwd
    DB_NAME=fintech
    DB_PORT=5442
    SRV_PORT=8000

## Up database

    docker-compose up

## Run the app

    go run cmd/app/main.go

----
# Rest API

----
Below you can read the descriptions of the endpoints calls, or you can see it here [Swagger](http://localhost:8000/swagger/index.html#/).

----
**Balance**
----
This option allow you to get user's balance by id.

* **URL**

  /users/balance

* **Method:**

  `POST`

*  **URL Params**

   None

* **Data Params**

  **Required:**
  ```
  {
  "id": 1
  }
  ```

* **Success Response:**

  If successful, then you should receive status code and response body.

    * **Code:** `200 OK`
    * **Content:** `{"id":,"amount":}`

* **Error Response:**

  In case of failure, you should receive status code and error message.

    * **Code:** `400 BAD REQUEST`
    * **Content:** `{"error": "grossbook get owner error: <user with this id doesn't exist>"}`

* **Sample Call:**

  ```
  curl --location --request POST 'localhost:8000/users/balance' \
    --header 'Content-Type: text/plain' \
    --data-raw '{
    "id": 200
    }'
  ```
  ----
**History**
----
This option allows you to get your transactions' history by id. You can limit transaction's quantity by "quantity" and you can sort output by "date" or "amount". 

* **URL**

  /users/history

* **Method:**

  `POST`

*  **URL Params**

   None

* **Data Params**

  **Required:**
  ```
  {
    "id": 2,
    "quantity": 1,
    "mode": "date"
  }
  ```

* **Success Response:**

  If successful, then you should receive only status code.

    * **Code:** `200 OK`
    * **Content:**
      ```
        [
          {
            "initiator_id": 2,
            "type": "WITHDRAW",
            "amount": 76.41201481050776,
            "timestamp": "2022-01-14T15:01:38.888762Z"
          }
        ]

* **Error Response:**

  In case of failure, you should receive status code and error message.

    * **Code:** `400 BAD REQUEST`
      **Content:** `{"error": "can't load history: <user with this id doesn't exist>"}`
    

* **Sample Call:**

  ```
  curl --location --request POST 'localhost:8000/users/history' \
    --header 'Content-Type: text/plain' \
    --data-raw '{
    "id": 200,
    "quantity": 1,
    "mode": "date"
    }'
  ```
  ----
**Deposit**
----
This option allows you to increase your balance by your id.

* **URL**

  /operations/deposit

* **Method:**

  `POST`

* **URL Params**

   None

* **Data Params**

  **Required:**
  ```
  {
    "initiator_id": 200,
    "amount": 100
  }
  ```

  * **Success Response:**

    If successful, then you should receive status code and response body.

      * **Code:** `201 CREATED`
      * **Content:**
          ```
            {
              "initiator": {
                "id": 200,
                "amount": 360.7355773908967
              },
              "type": "DEPOSIT",
              "amount": 100,
              "timestamp": "2022-01-14T16:10:52.3293451+03:00"
            }

* **Error Response:**

  In case of failure, you should receive status code and error message.

    * **Code:** `400 BAD REQUEST`
      **Content:** `{"error": "grossbook get user error: <user with this id doesn't exist>"}`

* **Sample Call:**

  ```
  curl --location --request POST 'localhost:8000/operations/deposit' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "initiator_id": 200,
  "amount": 100
  }'
  ```
  ----
**Withdraw**
----
This option allows you to decrease your balance by your id. You can choose currency through query.

* **URL**

  /operations/withdraw

* **Method:**

  `POST`

*  **URL Params**

   `?currency=USD`

* **Data Params**

  **Required:**
  ```
  {
    "initiator_id": 200,
    "amount": 1
  }
  ```

* **Success Response:**

  If successful, then you should receive status code and response body.

    * **Code:** `201 CREATED`
      * **Content:**
        ```
        {
          "initiator": {
            "id": 200,
            "amount": 360.7355773908967
          },
          "type": "WITHDRAW",
          "amount": 76.41201481050776,
          "timestamp": "2022-01-14T16:10:52.3293451+03:00"
        }

* **Error Response:**

  In case of failure, you should receive status code and error message.

    * **Code:** `400 BAD REQUEST`
      **Content:** `{"error": "grossbook get user error: <user with this id doesn't exist>"}`

* **Sample Call:**

  ```
  curl --location --request POST 'localhost:8000/operations/withdraw?currency=USD' \
  --header 'Content-Type: text/plain' \
  --data-raw '{
  "initiator_id": 200,
  "amount": 1
  }'
  ```

  ----
**Transfer**
----
This option allows you to transfer money from one user to another.

* **URL**

  /operations/transfer

* **Method:**

  `POST`

* **URL Params**

   None

* **Data Params**

  **Required:**
  ```
  {
    "initiator_id": 200,
    "receiver_id": 100,
    "amount": 1
  }
  ```

  * **Success Response:**

    If successful, then you should receive status code and response body.

      * **Code:** `201 CREATED`
      * **Content:**
      ```
        {
         "initiator": {
           "id": 200,
           "amount": 360.7355773908967
         },
         "type": "WITHDRAW",
         "amount": 76.41201481050776,
         "timestamp": "2022-01-14T16:10:52.3293451+03:00"
         }

* **Error Response:**

  In case of failure, you should receive status code and error message.

    * **Code:** `400 BAD REQUEST`
      **Content:** `{"error" :"grossbook get owner error: <user with this id doesn't exist>"}`

* **Sample Call:**

  ```
  curl --location --request POST 'localhost:8000/operations/transfer' \
  --header 'Content-Type: text/plain' \
  --data-raw '{
  "initiator_id": 200,
  "receiver_id": 100,
  "amount": 1
  }'
  ```
