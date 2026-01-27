package runner

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

// Runner handles the execution of the child command
type Runner struct {
	stderrBuf *CircularBuffer
}

// NewRunner creates a new Runner instance
func NewRunner() *Runner {
	return &Runner{
		stderrBuf: NewCircularBuffer(15), // Capture last 15 lines
	}
}

// Run executes the command and returns exit code, captured stderr, and error
func (r *Runner) Run(command []string) (int, string, error) {
	if len(command) == 0 {
		return 0, "", nil
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout // Stream directly to user terminal

	// Capture stderr
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return -1, "", err
	}

	// Tie stderr to both our buffer and os.Stderr (so user still sees it)
	// We use a MultiWriter to pipe to both.
	// However, we want to capture line by line for the buffer.
	// So we'll read from the pipe, write to os.Stderr, and feed our buffer.
	
	// Better approach:
	// We can't easily peek into os.Stderr if we just redirect it.
	// So we'll run a goroutine that copies from stderrPipe to os.Stderr AND r.stderrBuf
	
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// TeeReader reads from stderrPipe, writes to os.Stderr
		tee := io.TeeReader(stderrPipe, os.Stderr)
		// We copy from tee to our buffer
		// Note: CircularBuffer should implement io.Writer
		io.Copy(r.stderrBuf, tee)
	}()

	// Signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		if cmd.Process != nil {
			cmd.Process.Signal(sig)
		}
	}()

	if err := cmd.Start(); err != nil {
		return -1, "", err
	}

	wg.Wait() // Wait for stderr copying to finish
	err = cmd.Wait()

	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return exitCode, r.stderrBuf.String(), nil
}

// CircularBuffer to keep last N lines
type CircularBuffer struct {
	lines []string
	size  int
	buf   bytes.Buffer
}

func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		lines: make([]string, 0, size),
		size:  size,
	}
}

// Write implements io.Writer to capture output line by line
// This is a naive implementation; for full correctness dealing with partial writes
// we would need a proper line scanner. For simplicity here:
func (cb *CircularBuffer) Write(p []byte) (n int, err error) {
	// Write to internal buffer to handle partial lines? 
	// Or just append to lines?
	// Let's just accumulate bytes and split by newline for simplicity in this prototype
	
	n, err = cb.buf.Write(p)
	return n, err
}

// String returns the last N lines joined
func (cb *CircularBuffer) String() string {
	// Split whatever we have in buffer into lines
	// Note: this might be expensive for huge outputs but for a prototype it's fine.
	// Optimization: process lines as they come in.
	
    // Real implementation:
    // Just return the raw string if it's small, or tail it.
    // Given the requirement is just "Error Snippet", let's use a simple approach.
    // We'll just return the whole buffer if it's not huge, relying on the fact that
    // we typically want the ENTIRE stderr if it crashed, but the requirement said "Last 15 lines".
    
    // Let's do a simple line split here.
    lines := bytes.Split(cb.buf.Bytes(), []byte{'\n'})
    count := len(lines)
    if count == 0 {
        return ""
    }
    
    start := 0
    if count > cb.size {
        start = count - cb.size
    }
    
    return string(bytes.Join(lines[start:], []byte{'\n'}))
}
