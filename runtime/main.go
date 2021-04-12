package main

import (
	"runtime"
)

func main() {
	// Print the number of CPUs available
	println("Number of CPUs:", runtime.NumCPU())

	// Print the number of Goroutines
	println("Number of Goroutines:", runtime.NumGoroutine())

	// Print the Go version
	println("Go Version:", runtime.Version())

	// Print the operating system and architecture
	println("OS/Arch:", runtime.GOOS, "/", runtime.GOARCH)

	// Print the compiler used
	println("Compiler:", runtime.Compiler)

	// Print the number of cgo calls
	println("Number of cgo calls:", runtime.NumCgoCall())

	// Print the current memory statistics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	println("Alloc:", memStats.Alloc)
	println("TotalAlloc:", memStats.TotalAlloc)
	println("Sys:", memStats.Sys)
	println("NumGC:", memStats.NumGC)

	// Print the caller information
	pc, file, line, ok := runtime.Caller(0)
	if ok {
		fn := runtime.FuncForPC(pc)
		println("Caller Info:", fn.Name(), "at", file, ":", line)
	}

	// Print the stack trace of the current goroutine
	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, false)
	println("Stack Trace:\n", string(buf[:n]))

	// Force garbage collection
	println("Forcing garbage collection...")
	runtime.GC()

	// Yield the processor to allow other Goroutines to run
	println("Yielding processor...")
	runtime.Gosched()

	// Set the maximum number of CPUs to 1
	runtime.GOMAXPROCS(1)

	// Set the maximum number of CPUs that can be executing
	runtime.GOMAXPROCS(runtime.NumCPU())
}
