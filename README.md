
# mpay-fraud-service

This application is a bank transactions fraud detection service. Service has two seperate service endpoints and a Redis channel subscription.


## Endpoints
- /rules: Fraud detection parameters assignment endpoint
	- Request: 
	```json
	{
		"data": [
			{
				"name": "Rule 1",
				"key": "Key ",
				"priority": 1,
				"value": "10",
				"status": true,
				"operator": "gt"
			},{
				"name": "Rule 2",
				"key": "Key ",
				"priority": 4,
				"value": "10",
				"status": true,
				"operator": "gt"
			},{
				"name": "Rule 3",
				"key": "Key ",
				"priority": 3,
				"value": "10",
				"status": false,
				"operator": "gt"
			},{
				"name": "Rule 4",
				"key": "Key ",
				"priority": 2,
				"value": "10",
				"status": true,
				"operator": "gt"
			}
		]
	}
	```
- Response: Returns active rules sorted by priority ascanding
	```json
	{
		"Status": "success",
		"Code": 200,
		"Message": "Success",
		"Data": [
			{
				"name": "Rule 1",
				"key": "Key ",
				"priority": 1,
				"value": "10",
				"status": true,
				"operator": "gt"
			},
			{
				"name": "Rule 4",
				"key": "Key ",
				"priority": 2,
				"value": "10",
				"status": true,
				"operator": "gt"
			},
			{
				"name": "Rule 2",
				"key": "Key ",
				"priority": 4,
				"value": "10",
				"status": true,
				"operator": "gt"
			}
		]
	}
	```
  
| Name | Type | Value | Description |
|--|--|--|--|
| **Status** | string | success, fail, error  | enum value |
| **Code** | numeric | 200, 400, etc. | http status code |
| **Message** | string | message to the client  | Informal or exceptional messages |
| **Data** | string | rule list | active and sorted by priority rules list |

- /fraud: Fraud detection endpoint
	- Request: 
	```json
	{
		"data": {
			"amount": "250.00",
			"trx": "oaisufklafasfl1111112d1233",
			"card_number": "4943******2271",
			"expiration_month": "10",
			"expiration_year": "2022",
			"cardholders_name": "*** *****",
			"cvv": "402",
			"return_url": "https://envoysoft3.net/deposit/mpayReturn"
		},
		"user": {
			"username": "srdr16",
			"userID": "17206739184",
			"yearofbirth": "1983",
			"fullname": "*** *****",
			"email": "****@hotmail.com",
			"tckn": "17206739184",
			"ip_address": "95.10.24.238"
		}
	}
	```
- Response: Returns active rules sorted by priority ascanding
	```json
	{
		"Status": "success",
		"Code": 200,
		"Message": "Success",
		"Data": true
	}
	```
  
| Name | Type | Value | Description |
|--|--|--|--|
| **Status** | string | success, fail, error  | enum value |
| **Code** | numeric | 200, 400, etc. | http status code |
| **Message** | string | message to the client  | Informal or exceptional messages |
| **Data** | string | true, false | if passed from fraud detection then `true` else `false` |

## PUB/SUB
### Subscriptions
- fraud:rule_sets_changed :
	- Message: 
	```json
	{
		"data": [
			{
				"name": "Rule 1",
				"key": "Key ",
				"priority": 1,
				"value": "10",
				"status": true,
				"operator": "gt"
			},{
				"name": "Rule 2",
				"key": "Key ",
				"priority": 4,
				"value": "10",
				"status": true,
				"operator": "gt"
			},{
				"name": "Rule 3",
				"key": "Key ",
				"priority": 3,
				"value": "10",
				"status": false,
				"operator": "gt"
			},{
				"name": "Rule 4",
				"key": "Key ",
				"priority": 2,
				"value": "10",
				"status": true,
				"operator": "gt"
			}
		]
	}
	```
### Publishes
- fraud:blacklist : Used when it is determined that the relevant payload should be blacklisted
	- Message:
	```json
	```
- fraud:clear_counter : 
	- Message:
	```json
	```
- fraud:increase_counter : 
	- Message:
	```json
	```

## TO-DO
- ~~Integrate Redis~~
- ~~Create /rules http endpoint~~
- ~~Create /fraud http endpoint~~
- ~~Subscribe `fraud:rule_sets_changed` Redis channel~~ 
- ~~Publish events from appropriate Redis channel~~
- ~~Integrate generic rule engine~~
- Write rule-based fraud detection methods
- Write Unit Tests