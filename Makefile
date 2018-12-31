all:
	go build -o build/memo app.go controller.go db.go auth.go
