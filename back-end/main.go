package main

func main() {
	defer log_file.Close()
	logger.Info("Server started")
}
