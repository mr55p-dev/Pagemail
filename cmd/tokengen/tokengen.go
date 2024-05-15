package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"sync"
)

func main() {
	encoder := base64.NewEncoder(base64.URLEncoding, os.Stdout)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		io.Copy(encoder, os.Stdin)
		encoder.Close()
		wg.Done()
	}()
	wg.Wait()
	fmt.Fprintln(os.Stderr, "Done")
}
