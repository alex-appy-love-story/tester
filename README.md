# tester
This repository contains all of the test for the token shop system.

## Tests Overview

### Success Test
This test sends a valid request to the backend.
The payload of a successful request looks like this:
{
		Username: "Bob",
		TokenID:  4,
		Amount:   1,
}

### OutOfStock, InsufficientFunds, TokenNotFound Testa
These tests send requests that will create the different types of fails.

### Force Fail Microservices Test
Each request contains a TRIGGER_FAIL flag that will raise an error in one of the microservices specified by the flag. 

### Circuit Breaker Test
Sends 9 fail requests to observe the behavior of the circuit breaker.
The circuit breaker's open state is set to 8 seconds.
The circuit breaker opens after 5 consecutive fails.
