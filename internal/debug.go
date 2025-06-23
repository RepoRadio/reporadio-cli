package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

func isDebugEnabled() bool {
	return os.Getenv("DEBUG") != ""
}

func Debug(w io.Writer, msg string) {
	if !isDebugEnabled() {
		return
	}
	if w == nil {
		w = os.Stderr
	}
	fmt.Fprintf(w, "[DEBUG %s] %s\n", time.Now().Format("15:04:05"), msg)
}

func Debugf(w io.Writer, format string, args ...interface{}) {
	if !isDebugEnabled() {
		return
	}
	if w == nil {
		w = os.Stderr
	}
	fmt.Fprintf(w, "[DEBUG %s] %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, args...))
}

func DebugOpenAIRequest(w io.Writer, messages interface{}, model string) {
	if !isDebugEnabled() {
		return
	}
	if w == nil {
		w = os.Stderr
	}

	payload, _ := json.MarshalIndent(map[string]interface{}{
		"model":    model,
		"messages": messages,
	}, "", "  ")

	fmt.Fprintf(w, "[DEBUG %s] OpenAI Request:\n%s\n", time.Now().Format("15:04:05"), string(payload))
}

func DebugOpenAIResponse(w io.Writer, response interface{}) {
	if !isDebugEnabled() {
		return
	}
	if w == nil {
		w = os.Stderr
	}

	payload, _ := json.MarshalIndent(response, "", "  ")
	fmt.Fprintf(w, "[DEBUG %s] OpenAI Response:\n%s\n", time.Now().Format("15:04:05"), string(payload))
}
