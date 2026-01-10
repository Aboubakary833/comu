# Modular monolith following Clean architecture in Go - Learning project

This project is a learning-oriented modular monolith written in Go.
Its main goal is to explore clean modular boundaries, Clean architecture concepts, and scalable architecture patterns inside a single Go binary.


## ⚠️ Disclaimer

This project is not intended to be production-ready.
It is an experimental and educational project aimed at understanding how to structure a modular monolith in Go, in a context where real-world examples are rare.

Some architectural decisions may evolve as I advance.



## 🎯 Project goals

- Understand how to design a modular monolith

- Enforce clear boundaries between modules

- Practice dependency inversion without a framework

- Explore Clean architecture design

- Try to make modules communicate without event broker

- Prepare a possible future migration to microservices


## Folders structure

	- cmd/web
	- config
	- internal
		- modules
			- users
			- auth
			- posts
			- notifications
		- shared
	- migrations

**cmd/web**: This folder main function is the entry point of the whole application

**config**: Contains the configuration

**internal/modules**: That is the that contains the monolith modules

Each module follows the same structure:

- **domain**: Entities and business rules(interfaces)

- **application**: application logic(use cases)

- **infra**: repositories, external services

- **presentation**: HTTP handlers, routes, public API implementation

*module.go*: The module entry point and its public API interface and dependency exposure


**internal/shared**: Contains shared packages like logger(for now)
