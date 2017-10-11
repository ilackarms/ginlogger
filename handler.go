package ginlogger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/mattn/go-isatty"
)

// DisableConsoleColor disables color output in the console.
func DisableConsoleColor() {
	disableColor = true
}

func Handler(handler http.Handler, out io.Writer, notlogged ...string) http.Handler {
	isTerm := true

	if w, ok := out.(*os.File); !ok ||
		(os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd()))) ||
		disableColor {
		isTerm = false
	}

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Start timer
		start := time.Now()
		path := req.URL.Path
		raw := req.URL.RawQuery

		statusResponseWriter := &responseCodeInterceptor{ResponseWriter: w}

		// Process request
		handler.ServeHTTP(statusResponseWriter, req)

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)

			clientIP := clientIP(req)
			method := req.Method
			statusCode := statusResponseWriter.statusCode
			var statusColor, methodColor, resetColor string
			if isTerm {
				statusColor = colorForStatus(statusCode)
				methodColor = colorForMethod(method)
				resetColor = reset
			}
			//comment := c.Errors.ByType(ErrorTypePrivate).String()

			if raw != "" {
				path = path + "?" + raw
			}

			fmt.Fprintf(out, "[HTTP] %v |%s %3d %s| %13v | %15s |%s %-7s %s %s\n",
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, statusCode, resetColor,
				latency,
				clientIP,
				methodColor, method, resetColor,
				path,
				//comment,
			)
		}
	})
}
