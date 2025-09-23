# A Not-So-Basic Web Server

Software bloat often begins with architectural decisions. The software shelves at GitHub spill over with ready-made components of all kind, ready for the convenience shopper to grab. In no time, a software project that started out as a simple idea turns into a behemoth of services, infrastructure, and interdependencies.

And no one understands anymore what's going on.

## Go minimal

With this project, I want to go the opposite route: To build a platform with *just enough* features to serve as a web application platform or a web presence for a small business, and to actually understand what the code does. 

So here is a counterexample to the plethora of feature-loaded content management systems (CMSs): A super-minimal web server.

Well, it's not so super-minimal, as I added three features that the standard minimal web server examples omit. Two are crucial for any web server: proper timeouts and graceful shutdown. The third is a drastic simplification for the deployment. 

But let's dig into the code.

## A web server with an embedded file system

The web server shall serve just static HTML files, and nothing else. So I'll start with one of the aforementioned extra features: an embedded file system. The `embed` package lets you bake files and directories into the binary, for true single-binary deployment even if you need to deploy some non-Go files along with it. A `//go:` directive and the declaration of an embedded file system variable are all that's needed:

```go
//go:embed web/public/*
var fileSys embed.FS
```

Well, there's one more thing to do: In `func main`, I need to extract the `public` directory as a sub-filesystem; otherwise, the files would be served at `/web/public` rather than at the base URL `/`:

```go
	publicFS, err := fs.Sub(fileSys, "web/public")
```


Now I can set up a multiplexer and a handler that serves the sub-filesystem:

```go
	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServerFS(publicFS))
```

## Adding timeouts

Next step: create the HTTP server. I will not use the default server in `net/http`, as it lacks one of the crucial functionalities mentioned above: proper timeouts.

> Timeouts are possibly the most dangerous edge case to overlook.
>
> *–Filippo Valsorda*

Adding them is simple, though:

```go
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
```

([Source](https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/#:~:text=explicitly%20using%20a-,Server,-%3A))

Done! These three lines protect your server against slow clients and stalling connections. The settings are static and don't take varying payloads for different URL paths into account, but it's a start.

Now the server is ready to run:

```go
	go func() {
		log.Printf("Server starting on http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()
```

## Be nice to clients when shutting down

Finally, I'll add some logic for gracefully shutting down the server if the operating system sends a SIGTERM signal. A SIGTERM signal is sent to an app running in a terminal when you hit Ctrl+C, or to a background app through calling `kill <pid>` or `killall <process_name>`. 

The following code uses the `os/signal` package to catch a SIGTERM signal. Function `signal.Notify()` sends a value through a channel once it receives a signal. 

The code after the `Notify()` call waits on the channel before it calls the `Shutdown()` method from `net/http.Server` to close all listeners and idle connections and then wait *indefinitely* for active connections to return to idle. Well, indefinitely unless the context passed to it has a timeout. I add a 15-second timeout to be on the safe side for most connections (YYMV):

```go 
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutdown signal received, initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown due to error: %v", err)
	}
	log.Println("Server exited gracefully")
}
```

## Done!

This completes the first part of this lab. Try out the code: Clone the repository and call `go run .` in the root of the repository. Then open `http://localhost:8080` to verify that the server works as expected.

To test the file embedding, call `go build` and move the resulting binary to any directory of your liking, then invoke it from there. It'll work like it does inside the repo, as all required files in `web/public` are inside the binary. 

While this server is pretty bare-bones, it runs and serves HTML files! There are still a few things to add to make it ready for production tasks. Up next: a minimal UI.

## Dig deeper

Find out more about the packages used here:

- [embed package - embed](https://pkg.go.dev/embed)
- [http package - net/http](https://pkg.go.dev/net/http)
- [signal package - os/signal](https://pkg.go.dev/os/signal)