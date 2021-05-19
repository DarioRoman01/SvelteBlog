# SvelteBlog
![](https://github.com/Haizza1/my_first_repo/blob/master/svelte.png)
![](https://github.com/Haizza1/my_first_repo/blob/master/svelte2.png)
![](https://github.com/Haizza1/my_first_repo/blob/master/svelte3.png)

## What is this?
this can be call a personal project recreating a old project that i made in django,
the main difference between this project and the other project is that this project separates 
the frontend from the backend, the main reason for this is learn how to build a fullstack project
building my own backend and frontend, work with pagination in both sizes, sessions and a lot more.

## Project status
This project can be considerer finished

## Tecnologies used for this project

### Backend
* Golang
* PostgresSQL
* Echo
* Gorm
* JWT

### Frontend
* TypeScript
* Svelte
* Routify
* Svelte Materialify

## Usage

First clone the repo
```
$ git clone https://github.com/Haizza1/SvelteBlog.git
```

### Backend setup
```
$ cd server

$ go mod download
``` 

before you start the server check the .env.example file and fill the all the fields

you will notice this 
```
EMAIL_KEY=
```
if you want to try locally to register a user you will ned access to mailslurp api, but dont worry its free! check https://www.mailslurp.com/ then you have register and get the api key, then put that api key in the env variable

then just run 
```
$ go run main.go
```

### Frontend setup
```
$ cd web

$ yarn

$ yarn dev
```
