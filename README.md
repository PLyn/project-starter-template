# templ-quickstart

## Introduction

Custom starter template to get up and running. The underlying stack is Go, Templ, Datastar and tailwind.

## Core Technologies

As mentioned above, this project depends on some awesome technologies. Let me start by giving credit where credit is due:

- [Go](https://go.dev/) - Version 1.23 or greater required
- [Templ](https://templ.guide/)
- [Datastar](https://data-star.dev/)
- [Tailwindcss](https://tailwindcss.com/)

## Installation

### Clone the Repository

```bash
git clone https://github.com/PLyn/project-starter-template <target-directory>
```

```bash
cd <target-directory>
```

### Install Dependencies

```bash
go mod tidy
```

### Create a .env file and include a PORT variable

```bash
touch .env; 
```

```bash
echo "PORT=8080" > .env
```

## Build Steps and Serving

This project requires a build step. The following are commands needed to build your html and css output.

### Templ HTML Generation

With templ installed and the binary somewhere on your PATH, run the following to generate your HTML components and templates (remove --watch to simply build and not hot reload)

```bash
templ generate --watch
```

### CSS File Generation

With the [Tailwind Binary](https://tailwindcss.com/blog/standalone-cli) installed and moved somewhere on your PATH, run the following to generate your CSS output for your tailwind classes (remove --watch to simply build and not hot reload)

```bash
tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch
```

### Serving with Air

With the [Air Binary](https://github.com/cosmtrek/air) installed and moved somewhere on your PATH, run the following to serve and hot reload the application:

```bash
air
```

To configure air, you can modify .air.toml in the root of the project. (it will be auto-generated after the first time you run air in your repo)

## Project Overview

This project has a few core concepts to help you get going, let's start with ./main.go

### Main - ./main.go

This is our applications entry-point and does a few things:

1. Here, we load in our .env file and then we initialize our mux server.

```go
_ = godotenv.Load()
mux := http.NewServeMux()
```

2. We define a few basic routes for our server. I will go into these routes in more depth later. In short, these routes enable you to use static files in your project, to use a favicon.ico, and sets up a view found at "/".

```go
mux.HandleFunc("GET /favicon.ico", view.ServeFavicon)
mux.HandleFunc("GET /static/", view.ServeStaticFiles)
mux.Handle("/", middleware.AdaptHandler(handlers.MyHandler))
```

3. We serve our application on the PORT defined at ./.env

```go
fmt.Println(fmt.Sprintf("server is running on port %s", os.Getenv("PORT")))
err := http.ListenAndServe(":"+os.Getenv("PORT"), mux)
if err != nil {
    fmt.Println(err)
}
```
### Middleware - ./internal/middleware/middleware.go

Custom middleware can be implemented with ease in this project. There is a custom type that has the signature: 

```go
type CustomHandler func(ctx *CustomContext, w http.ResponseWriter, r *http.Request)
```
which means we can create each of our handlers for the application like this:

```go
// Example function using CustomHandler type.
func MyHandler(ctx *middleware.CustomContext, w http.ResponseWriter, r *http.Request) {
	...Handler logic here
}
```

The reason for doing this is to be able to add a new endpoint in a really simple way. Here is what a possible endpoint with 2 middleware functions to run:

```go
mux.Handle("/", middleware.AdaptHandler(handlers.MyHandler, AuthMiddleware, LoggingMiddleware))
```

This is all possible because of the middleware function called AdaptHandler. Here is the main piece that ties it together:

```go
func AdaptHandler(handler CustomHandler, middleware ...CustomMiddleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customContext := &CustomContext{
			Context:   context.Background(),
			StartTime: time.Now(),
			...Any other fields needed throughout your app
		}

		//Run all custom middleware functions before Handler runs
		for _, mw := range middleware {
			err := mw(customContext, w, r)
			if err != nil {
				return
			}
		}
		//run Handler logic
		handler(customContext, w, r)

		//Run logic after Handler runs
		Log(customContext, w, r)
	}
}
```

If a middleware returns an error, the chain will stop executing. This enables you to allow your middleware to write responses early and avoid calling the handler in case of an error.

### Creating Custom Middleware

Let's say you want to create custom middleware. Here is how to do so:

1. If this middleware requires some context, add the context value to the CustomContext type.

```go
type CustomContext struct {
    context.Context
    StartTime time.Time
    NewContextValue string
}
```

2. Define your new middleware functions (remember middleware must match the CustomMiddleware type definition).

```go
// this middleware will be placed early in the chain
func EarlyMiddleware(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error {
	ctx.NewContextValue = "I was set early in the chain" // set your new context value
	return nil
}

// this middleware will be place late in the chain
func LateMiddleware(ctx *CustomContext, w http.ResponseWriter, r *http.Request) error {
	fmt.Println(ctx.NewContextValue) // outputs "I was set early in the chain"
	return nil
}
```

3. Include the middleware in your Chain func in your routes.

```go
// modified version of ./main.go
mux.Handle("/", middleware.AdaptHandler(handlers.MyHandler, middleware.EarlyMiddleware, middleware.LateMiddleware))
})
```

That's it! Easily create custom middleware without the need to deeply nest your routes.

### SSE endpoints - ./internal/SSE/*.go

Datastar heavily uses SSE endpoints and you create a new one as needed to stream raw HTML/content/data/notifications/etc to your web page in realtime.

```go
func Backup(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	var buf bytes.Buffer
	for {
		buf.Reset()
		id := uuid.New()
		replaceHTML := components.Replace("feed", id.String(), time.Now().UTC())
		replaceHTML.Render(context.Background(), &buf)

		sse.MergeFragments(fmt.Sprint(buf.String()),)
		time.Sleep(1 * time.Second)
	}
}

```

### Templates - ./internal/template/

The template folder is used for subdirectories for each endpoint to localize all logic and components for that endpoint in a directory. All shared components are placed in the sharedComponents directory if they are used across endpoints. This is done using Templ templating language/

The base.templ in the sharedComponents is the Base webpage with the header and footer of every page and any required HTML/scrips/CSS is loaded or written here.

```html
templ Base(title string) {
	<html>
		<head>
			<meta charset="UTF-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
			<link rel="stylesheet" href="/static/css/output.css"/>
			<title>{ title }</title>
			<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.0-beta.1/bundles/datastar.js"></script>
		</head>
		@Banner()
		<body>
			<main class="p-6 grid gap-4">
				{ children... }
			</main>
		</body>
	</html>
}
```

Also note, datastar and your tailwind output are included in the head of this template:

```html
<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.0-beta.1/bundles/datastar.js"></script>
<link rel="stylesheet" href="/static/css/output.css"></link>
```
