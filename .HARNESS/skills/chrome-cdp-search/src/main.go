package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

type Options struct {
	Query, Engine, Browser string
	Count                  int
	Timeout                time.Duration
}
type Result struct {
	Rank    int    `json:"rank"`
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}
type Output struct {
	Status       string   `json:"status"`
	Engine       string   `json:"engine"`
	Query        string   `json:"query"`
	SearchEngine string   `json:"searchEngine"`
	ResultsCount int      `json:"resultsCount"`
	Results      []Result `json:"results"`
	Error        string   `json:"error,omitempty"`
}
type cdpResponse struct {
	ID     uint64          `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Message string `json:"message"`
	} `json:"error"`
}
type targetInfo struct {
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

func parseArgs(args []string) (Options, error) {
	fs := flag.NewFlagSet("chrome-search", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var o Options
	var timeout string
	fs.StringVar(&o.Query, "query", "", "search query")
	fs.StringVar(&o.Query, "q", "", "search query")
	fs.StringVar(&o.Engine, "engine", "bing", "bing, google, or duckduckgo")
	fs.IntVar(&o.Count, "count", 10, "maximum results (1-50)")
	fs.StringVar(&o.Browser, "browser", "", "browser executable (default: discover system Chromium/Chrome)")
	fs.StringVar(&timeout, "timeout", "30s", "overall timeout (for example 30s)")
	if err := fs.Parse(args); err != nil {
		return o, err
	}
	if strings.TrimSpace(o.Query) == "" {
		return o, errors.New("--query is required")
	}
	switch o.Engine {
	case "bing", "google", "duckduckgo":
	default:
		return o, fmt.Errorf("unsupported engine %q (choose bing, google, duckduckgo)", o.Engine)
	}
	if o.Count < 1 || o.Count > 50 {
		return o, errors.New("--count must be between 1 and 50")
	}
	d, err := time.ParseDuration(timeout)
	if err != nil || d <= 0 {
		return o, errors.New("--timeout must be a positive duration")
	}
	o.Timeout = d
	return o, nil
}
func searchURL(engine, query string) string {
	q := url.QueryEscape(query)
	switch engine {
	case "google":
		return "https://www.google.com/search?q=" + q + "&num=50"
	case "duckduckgo":
		return "https://html.duckduckgo.com/html/?q=" + q
	default:
		return "https://www.bing.com/search?q=" + q + "&count=50"
	}
}
func browserPath(requested string) (string, error) {
	if requested != "" {
		p, e := exec.LookPath(requested)
		if e == nil && os.PathSeparator == '\\' && strings.EqualFold(filepath.Base(requested), filepath.Base(p)) && filepath.Ext(requested) != filepath.Ext(p) {
			p = requested
		}
		if e != nil {
			return "", fmt.Errorf("browser %q not found: %w", requested, e)
		}
		return p, nil
	}
	probe := []string{"chromium", "chromium-browser", "google-chrome", "google-chrome-stable", "chrome", "msedge"}
	if os.PathSeparator == '\\' {
		for i, n := range probe {
			if filepath.Ext(n) == "" {
				probe[i] = n + ".exe"
			}
		}
	}
	for _, name := range probe {
		if p, e := exec.LookPath(name); e == nil {
			return p, nil
		}
	}
	return "", errors.New("Chromium/Chrome not found; set --browser")
}
func main() {
	o, err := parseArgs(os.Args[1:])
	if err != nil {
		emit(Output{Status: "error", Error: err.Error()})
		os.Exit(2)
	}
	base, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()
	ctx, cancel := context.WithTimeout(base, o.Timeout)
	defer cancel()
	results, err := runSearch(ctx, o)
	if err != nil {
		emit(Output{Status: "error", Engine: "Chromium/Chrome DevTools Protocol", Query: o.Query, SearchEngine: o.Engine, Error: err.Error()})
		os.Exit(1)
	}
	emit(Output{Status: "success", Engine: "Chromium/Chrome DevTools Protocol", Query: o.Query, SearchEngine: o.Engine, ResultsCount: len(results), Results: results})
}
func emit(v Output) { enc := json.NewEncoder(os.Stdout); enc.SetIndent("", "  "); _ = enc.Encode(v) }
func runSearch(ctx context.Context, o Options) ([]Result, error) {
	browser, err := browserPath(o.Browser)
	if err != nil {
		return nil, err
	}
	profile, err := os.MkdirTemp("", "chrome-search-profile-")
	if err != nil {
		return nil, fmt.Errorf("create temporary profile: %w", err)
	}
	if wd, e := os.Getwd(); e == nil {
		if fi, e2 := os.Stat(profile); e2 != nil || !fi.IsDir() {
			return nil, fmt.Errorf("temporary profile %q is not usable from the child process (working directory %q, os.TempDir=%q, TEMP=%q): export a writable TMPDIR/TEMP and retry: %v", profile, wd, os.TempDir(), os.Getenv("TEMP"), e2)
		}
	}
	defer os.RemoveAll(profile)
	cmd := exec.Command(browser, "--headless=new", "--disable-gpu", "--no-first-run", "--no-default-browser-check", "--disable-extensions", "--remote-debugging-port=0", "--user-data-dir="+profile, "about:blank")
	if err = cmd.Start(); err != nil {
		return nil, fmt.Errorf("start browser: %w", err)
	}
	defer func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }()
	stop := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = cmd.Process.Kill()
		case <-stop:
		}
	}()
	defer close(stop)
	portFile := filepath.Join(profile, "DevToolsActivePort")
	var port string
	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		b, e := os.ReadFile(portFile)
		if e == nil {
			port = strings.TrimSpace(strings.Split(string(b), "\n")[0])
			if port != "" {
				break
			}
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("browser startup timed out: %w", ctx.Err())
		case <-time.After(50 * time.Millisecond):
		}
	}
	client := &http.Client{Timeout: 2 * time.Second}
	var target targetInfo
	for {
		req, _ := http.NewRequestWithContext(ctx, http.MethodPut, "http://127.0.0.1:"+port+"/json/new?about:blank", nil)
		resp, e := client.Do(req)
		if e == nil {
			_ = json.NewDecoder(resp.Body).Decode(&target)
			_ = resp.Body.Close()
			if target.WebSocketDebuggerURL != "" {
				break
			}
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("create DevTools page: %w", ctx.Err())
		case <-time.After(50 * time.Millisecond):
		}
	}
	ws, err := dialWebSocket(ctx, target.WebSocketDebuggerURL)
	if err != nil {
		return nil, err
	}
	defer ws.Close()
	var id uint64
	call := func(method string, params any) (json.RawMessage, error) {
		n := atomic.AddUint64(&id, 1)
		b, _ := json.Marshal(map[string]any{"id": n, "method": method, "params": params})
		if err := ws.write(b); err != nil {
			return nil, err
		}
		for {
			raw, e := ws.read(ctx)
			if e != nil {
				return nil, e
			}
			var r cdpResponse
			if json.Unmarshal(raw, &r) != nil || r.ID != n {
				continue
			}
			if r.Error != nil {
				return nil, errors.New(r.Error.Message)
			}
			return r.Result, nil
		}
	}
	if _, err = call("Page.enable", map[string]any{}); err != nil {
		return nil, err
	}
	if _, err = call("Runtime.enable", map[string]any{}); err != nil {
		return nil, err
	}
	if _, err = call("Page.navigate", map[string]any{"url": searchURL(o.Engine, o.Query)}); err != nil {
		return nil, err
	}
	// Poll DOM until at least one organic result exists; evaluation runs in Chromium so real DOM parsing needs no Go HTML dependency.
	deadline := time.NewTimer(8 * time.Second)
	defer deadline.Stop()
	tick := time.NewTicker(200 * time.Millisecond)
	defer tick.Stop()
	for {
		expr := resultExpression(o.Engine, o.Count)
		raw, e := call("Runtime.evaluate", map[string]any{"expression": expr, "returnByValue": true})
		if e == nil {
			var ev struct {
				Result struct {
					Value json.RawMessage `json:"value"`
				} `json:"result"`
			}
			if json.Unmarshal(raw, &ev) == nil && len(ev.Result.Value) > 0 {
				var rs []Result
				if json.Unmarshal(ev.Result.Value, &rs) == nil && len(rs) > 0 {
					return rs, nil
				}
			}
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("search timed out: %w", ctx.Err())
		case <-deadline.C:
			return nil, errors.New("no organic results extracted (blocked page, consent wall, or changed search markup)")
		case <-tick.C:
		}
	}
}
func resultExpression(engine string, count int) string {
	selectors := map[string]string{"bing": "li.b_algo", "google": "div.MjjYud, div.g", "duckduckgo": ".result"}
	s, _ := json.Marshal(selectors[engine])
	return fmt.Sprintf(`(()=>{const sel=%s, limit=%d, out=[]; for(const el of document.querySelectorAll(sel)){if(out.length>=limit)break;let a,t,sn;if(%q==='bing'){a=el.querySelector('h2 a');t=el.querySelector('h2');sn=el.querySelector('.b_caption p,.b_caption,.b_snippet')}else if(%q==='google'){a=el.querySelector('a:has(h3)');t=el.querySelector('h3');sn=el.querySelector('.VwiC3b,.yXK7lf')}else{a=el.querySelector('a.result__a');t=a;sn=el.querySelector('.result__snippet')}if(!a||!t||!a.href||!/^https?:/.test(a.href))continue;out.push({rank:out.length+1,title:(t.innerText||t.textContent||'').trim(),url:a.href,snippet:sn?(sn.innerText||sn.textContent||'').trim():''})}return out})()`, s, count, engine, engine)
}

type websocket struct {
	c net.Conn
	r *bufio.Reader
}

func dialWebSocket(ctx context.Context, raw string) (*websocket, error) {
	u, e := url.Parse(raw)
	if e != nil {
		return nil, e
	}
	host := u.Host
	d := net.Dialer{}
	c, e := d.DialContext(ctx, "tcp", host)
	if e != nil {
		return nil, fmt.Errorf("connect DevTools websocket: %w", e)
	}
	keyb := make([]byte, 16)
	_, _ = rand.Read(keyb)
	key := base64.StdEncoding.EncodeToString(keyb)
	path := u.RequestURI()
	if path == "" {
		path = "/"
	}
	_, e = fmt.Fprintf(c, "GET %s HTTP/1.1\r\nHost: %s\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: %s\r\nSec-WebSocket-Version: 13\r\n\r\n", path, host, key)
	if e != nil {
		c.Close()
		return nil, e
	}
	r := bufio.NewReader(c)
	resp, e := http.ReadResponse(r, &http.Request{Method: "GET"})
	if e != nil {
		c.Close()
		return nil, e
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		c.Close()
		return nil, fmt.Errorf("DevTools websocket upgrade failed: %s", resp.Status)
	} // Validate the server accept hash from the request key.
	h := sha1.Sum([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	if resp.Header.Get("Sec-WebSocket-Accept") != base64.StdEncoding.EncodeToString(h[:]) {
		c.Close()
		return nil, errors.New("invalid DevTools websocket accept")
	}
	return &websocket{c: c, r: r}, nil
}
func (w *websocket) Close() error { return w.c.Close() }
func (w *websocket) write(p []byte) error {
	mask := make([]byte, 4)
	if _, e := rand.Read(mask); e != nil {
		return e
	}
	h := []byte{0x81}
	n := len(p)
	if n < 126 {
		h = append(h, 0x80|byte(n))
	} else if n <= 65535 {
		h = append(h, 0x80|126, byte(n>>8), byte(n))
	} else {
		return errors.New("CDP message too large")
	}
	h = append(h, mask...)
	for i := range p {
		p[i] ^= mask[i%4]
	}
	if _, e := w.c.Write(append(h, p...)); e != nil {
		return e
	}
	return nil
}
func (w *websocket) read(ctx context.Context) ([]byte, error) {
	_ = w.c.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	for {
		h := make([]byte, 2)
		if _, e := io.ReadFull(w.r, h); e != nil {
			if ne, ok := e.(net.Error); ok && ne.Timeout() {
				if ctx.Err() != nil {
					return nil, ctx.Err()
				}
				continue
			}
			return nil, e
		}
		n := uint64(h[1] & 0x7f)
		if n == 126 {
			b := make([]byte, 2)
			if _, e := io.ReadFull(w.r, b); e != nil {
				return nil, e
			}
			n = uint64(b[0])<<8 | uint64(b[1])
		} else if n == 127 {
			b := make([]byte, 8)
			if _, e := io.ReadFull(w.r, b); e != nil {
				return nil, e
			}
			for _, x := range b {
				n = n<<8 | uint64(x)
			}
		}
		if n > 16<<20 {
			return nil, errors.New("CDP frame too large")
		}
		if h[1]&0x80 != 0 {
			return nil, errors.New("masked server websocket frame")
		}
		p := make([]byte, n)
		if _, e := io.ReadFull(w.r, p); e != nil {
			return nil, e
		}
		switch h[0] & 0x0f {
		case 1:
			return p, nil
		case 8:
			return nil, io.EOF
		case 9:
			continue
		default:
			continue
		}
	}
}
