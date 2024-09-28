package appstates

import (
	"fmt"
	"os"
)

const (
	failureTag string = "[ERROR]"
	successTag string = "[DONE]"
)

func printFailedAndPanic(code int, message string, details string) {
	fmt.Printf("%s | %s | %d | %s\n", failureTag, message, code, details)
	os.Exit(code)
}

func printSuccess(code int, message string, details string) {
	fmt.Printf("%s | %s | %d | %s\n", successTag, message, code, details)
}

// Unknown and unexpected reason crashed the app
func PanicUnexpected(details string) {
	printFailedAndPanic(1, "unexpected reason", details)
}

// Environment or config parameters are missing
func PanicMissingEnvParams(details string) {
	printFailedAndPanic(400, "missing environment params", details)
}

// Server socket error
func PanicServerSocketFailure(details string) {
	printFailedAndPanic(500, "server socket failure", details)
}

// Server internal error
func PanicServerInternalError(details string) {
	printFailedAndPanic(501, "server internal error", details)
}

// DB connection failed
func PanicDBConnectionFailed(details string) {
	printFailedAndPanic(600, "db connection failed", details)
}

// Config parameters extracted
func DoneConfigExtraction(details string) {
	printSuccess(100, "config extracted", details)
}

// Database connected
func DoneDbConnection(details string) {
	printSuccess(200, "db connected", details)
}

// Server launched
func DoneServerLaunch(details string) {
	printSuccess(300, "server launched", details)
}
