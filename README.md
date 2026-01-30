# Modular monolith following Clean architecture in Go - Learning project

This project is a learning-oriented modular monolith written in Go.
Its main goal is to explore clean modular boundaries, Clean architecture concepts, and scalable architecture patterns.


## ⚠️ Disclaimer

This project is not intended to be production-ready.
It is an experimental and educational project aimed at understanding how to structure a modular monolith in Go, in a context where real-world examples are rare.


## 🎯 Project goals

- Understand how to design a modular monolith

- Enforce clear boundaries between modules

- Explore Clean architecture design

- Make modules communicate through public apis

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

**internal/modules**: That is the foloder that contains the monolith modules

Each module follows the same structure:

- **domain**: Entities and business rules(interfaces)

- **application**: application logic(use cases)

- **infra**: repositories, external services

- **presentation**: HTTP handlers, routes...(UI also should goes there, but hey this is a simple Rest API)

*module.go*: Module entry point
*api.go*: Module public API interface and dependency exposure


**internal/shared**: Contains shared packages like logger, validator and utils


## Routes

**Login**:

	POST 	/login
	POST 	/login/verify
	POST 	/login/resend_otp
	POST 	/login/refresh

**Register**:

	POST 	/register
	POST 	/register/verify
	POST 	/register/resend_otp

**Reset Password**:

	POST 	/reset_password
	POST 	/reset_password/verify
	POST 	/reset_password/resend_otp
	POST 	/reset_password/new_password

**Posts**:

	GET 	/posts
	POST 	/posts/create
	GET 	/posts/:slug
	PUT 	/posts/update/:post_id
	DELETE  /posts/delete/:post_id

**Comments**:

	GET		/comments/list/post_id
	POST	/comments/create
	POST	/comments/update/:comment_id
	DELETE	/comments/delete/:comment_id



## Running the project

```sh
	mv .env.example .env && docker compose up -d
