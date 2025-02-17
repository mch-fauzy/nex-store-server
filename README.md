# Nexmedis BE Technical Test

# 1

## Question
You are tasked with designing an API for an e-commerce platform. The system must support the following features:
- User registration and authentication
- Viewing and searching products
- Adding items to a shopping cart
- Completing a purchase

Design the RESTful endpoints for the above features. Describe your choice of HTTP methods (GET, POST, PUT, DELETE), URL structure, and the expected response formats. Assume that users need to authenticate before performing certain actions (e.g., adding items to the cart)

## Answer

### Table of Contents

- [Technologies Used](#technologies-used)
- [Setup](#setup)
- [API Collections](#API-collections)

### Technologies Used
- Go version 1.22.x
- Docker version 4.28
- Postmann version 10.24.x

### Setup

1. Clone this repository:

   ```
   git clone https://github.com/mch-fauzy/nexmedis-be-technical-test.git
   ```

2. Navigate to the project directory:

   ```
   cd nexmedis-be-technical-test
   ```

3. To start the application, run the following command in the project root folder:

   ```
   docker-compose up --build
   ```

### API Collections

To simplify testing of the API endpoints, a Postman collections is provided. Follow the steps below to import and use it:

1. Use the Postman collection JSON file [**nexmedis-be-technical-test.postman_collection.json**](docs/nexmedis-be-technical-test.postman_collection.json) in this project directory

2. Open Postman

3. Click on the "Import" button located at the top left corner of the Postman interface

4. Select the JSON file

5. Once imported, you will see a new collection named "nexmedis-be-technical-test" in your Postman collections

6. You can now use this collection to test the API endpoints by sending requests to your running API server

# 2 

## Question
Consider a database table Users with the following columns:
- id (Primary Key) 
- username 
- email 
- created_at 

Your task is to design an indexing strategy to optimize the following queries:
- Fetch a user by username
- Fetch users who signed up after a certain date (created_at > "2023-01-01")
- Fetch a user by email

Explain which columns you would index, and whether composite indexes or individual indexes would be appropriate for each query. Discuss trade-offs in terms of read and write performance

## Answer

```
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT,
    email TEXT,
    created_at TIMESTAMPTZ(3) DEFAULT now()
);

-- Create unique index on username and email
CREATE UNIQUE INDEX idx_users_username ON users (username);
CREATE UNIQUE INDEX idx_users_email ON users (email);

-- Create an index on the created_at column
CREATE INDEX idx_users_created_at ON users (created_at);
```

- username and email:
Since both queries "Fetch a user by username" and "Fetch a user by email" require searching by a single column, and if these fields are unique, then adding unique constraints will creates unique indexes. This ensures fast lookups

- created_at:
For the query "Fetch users who signed up after a certain date", create an index on the created_at column to speed up the range condition

In this scenario, each query filters by a single column. Composite indexes are most beneficial when multiple columns are queried together. Adding unnecessary composite indexes can increase write overhead without a proportional benefit in read performance

Indexes significantly improve query speed for lookups and range queries, while this speeds up read operations, every insert, update, or delete  must update the index, which can slow down writes. However, for many applications where read performance is critical, this is an acceptable trade-off

# 3 

## Question
You need to implement a function that simulates a bank account system. Multiple users can simultaneously access and update their account balance. Your system must ensure that concurrent access does not result in race conditions. Implement the function that:
- a. Deposits money into an account.
- b. Withdraws money from an account (ensuring there’s enough balance).

## Answer

You can see the implementation on [**transaction-service.go**](services/transaction-service.go)

- a. Deposit money:

```
package services

import (
	"github.com/nexmedis-be-technical-test/models"
	"github.com/nexmedis-be-technical-test/models/dto"
	"github.com/rs/zerolog/log"
)

func (s *Service) TransactionTopUpBalanceByUserId(req dto.TransactionTopUpBalanceByUserIdRequest) (dto.TransactionTopUpBalanceByUserIdResponse, error) {
	var response dto.TransactionTopUpBalanceByUserIdResponse
	users, _, err := s.Repository.UserFindManyAndCountByFilter(models.Filter{
		FilterFields: []models.FilterField{
			{
				Field:    models.UserDbField.Id,
				Operator: models.OperatorEqual,
				Value:    req.UserId,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionTopUpBalanceByUserId] Service error getting users")
		return response, err
	}

	err = s.Repository.UserUpdateById(models.UserPrimaryId{Id: users[0].Id}, &models.User{
		Balance: users[0].Balance + float32(req.TopupAmount),
	})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionTopUpBalanceByUserId] Service error updating user balance")
		return response, err
	}

	user, err := s.Repository.UserFindById(models.UserPrimaryId{Id: users[0].Id})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionTopUpBalanceByUserId] Service error retrieving user by id")
		return response, err
	}

	response = dto.TransactionTopUpBalanceByUserIdResponse{
		Balance: user.Balance,
	}

	return response, nil
}

```

- b. Withdraw balance

```
package services

import (
	"github.com/nexmedis-be-technical-test/models"
	"github.com/nexmedis-be-technical-test/models/dto"
	"github.com/nexmedis-be-technical-test/utils/failure"
	"github.com/rs/zerolog/log"
)

func (s *Service) TransactionWithdrawBalanceByUserId(req dto.TransactionWithdrawBalanceByUserIdRequest) (dto.TransactionWithdrawBalanceByUserIdResponse, error) {
	var response dto.TransactionWithdrawBalanceByUserIdResponse

	users, _, err := s.Repository.UserFindManyAndCountByFilter(models.Filter{
		FilterFields: []models.FilterField{
			{
				Field:    models.UserDbField.Id,
				Operator: models.OperatorEqual,
				Value:    req.UserId,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionWithdrawBalanceByUserId] Service error getting user")
		return response, err
	}

	currentBalance := users[0].Balance
	withdrawAmount := float32(req.WithdrawAmount)
	if currentBalance < withdrawAmount {
		err = failure.BadRequest("Insufficient balance")
		return response, err
	}

	err = s.Repository.UserUpdateById(models.UserPrimaryId{Id: users[0].Id}, &models.User{
		Balance:   currentBalance - withdrawAmount,
		UpdatedBy: req.Email,
	})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionWithdrawBalanceByUserId] Service error updating user balance")
		return response, err
	}

	user, err := s.Repository.UserFindById(models.UserPrimaryId{Id: users[0].Id})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionWithdrawBalanceByUserId] Service error retrieving user by id")
		return response, err
	}

	response = dto.TransactionWithdrawBalanceByUserIdResponse{
		Balance: user.Balance,
	}

	return response, nil
}

```

# 4 

## Question
Given a database table orders with the following schema:
```
CREATE TABLE orders (
    id INT PRIMARY KEY,
    customer_id INT,
    product_id INT,
    order_date TIMESTAMP,
    amount DECIMAL(10, 2)
);

*Assume that customer_id is indexed, but amount and order_date are not indexed
```


- a. Write an optimized SQL query to find the top 5 customers who spent the most money in the past month. 
- b. How would you improve the performance of this query in a production environment? 

## Answer

- a. SQL Query:

```
SELECT 
    customer_id, 
    SUM(amount) AS total_spent
FROM orders
WHERE order_date >= NOW() - INTERVAL '1 month'
GROUP BY customer_id
ORDER BY total_spent DESC
LIMIT 5;
```

- b. Improve performance:
  - Create an index on order_date because the query filters rows based on the that column: `CREATE INDEX idx_orders_order_date ON orders(order_date);`
  - Create a composite index because the query uses order_date for filtering and customer_id for grouping while reading amount for aggregation. A composite index that includes these columns can speed up data retrieval: `CREATE INDEX idx_orders_date_customer_amount ON orders(order_date, customer_id, amount);`
  - Use pre-aggregate data if the query is run frequently by create job in database level and save the result in different table

# 5 

## Question
You are tasked with refactoring a monolithic service that handles multiple responsibilities such as authentication, file uploads, and user data processing. The system has become slow and hard to maintain. How would you approach refactoring the service?

- a. What steps would you take to decompose the service into smaller, more manageable services?
- b. How would you ensure that the new system is backward compatible with the old one during the transition?

## Answer

- a. Steps:
	1. Break the service into logical modules based on functionality (e.g., authentication, file uploads, user data processing). For instance, authentication might have its own rules and data, while file uploads involve storage and processing, and user data processing might have its own business logic
	2. Organize the code into layers (e.g., service, model, handlers, infrastructure, etc). This helps isolate concerns and makes the code easier to understand and maintain
	3. Establish clear, well-documented APIs for each service
	4. Gradually replace old code without a full rewrite at once
	5. Deploy the new services to a small subset of users first and monitor for issues

- b. Ensure backward compatibility:
	- Add versioned endpoints (e.g., /v1/auth vs. /v2/auth) so that existing system can continue using the old endpoints until they are ready to fully migrate
	- Introduce feature toggles to control which parts of the system are handled by new services. This allows you to enable or disable the new functionality
	- Maintain rollback strategies if unforeseen issues occur during the transition
	- Update the documentation to aid in troubleshooting and future maintenance