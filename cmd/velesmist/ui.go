// SPDX-FileCopyrightText: 2026 VelesMist contributors
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/asketmc/VelesMist/internal/config"
	apperrors "github.com/asketmc/VelesMist/internal/errors"
	"github.com/asketmc/VelesMist/internal/report"
)

const defaultUIAddr = "127.0.0.1:8765"

type uiConfig struct {
	Addr string
	Open bool
}

type uiScanFunc func(context.Context, config.ScanConfig) (report.ScanResult, error)

type uiScanRequest struct {
	SteamID    string `json:"steam_id"`
	Game       string `json:"game"`
	Fixture    string `json:"fixture"`
	PriceCache string `json:"price_cache"`
	MinPrice   string `json:"min_price"`
	Currency   string `json:"currency"`
	NoCache    bool   `json:"no_cache"`
}

type uiErrorResponse struct {
	OK       bool   `json:"ok"`
	Kind     string `json:"kind"`
	ExitCode int    `json:"exit_code"`
	Message  string `json:"message"`
}

type uiScanResponse struct {
	OK     bool              `json:"ok"`
	Report report.ScanResult `json:"report"`
}

func runUI(args []string, stdout io.Writer, stderr io.Writer) int {
	cfg, err := parseUI(args)
	if err != nil {
		fmt.Fprintf(stderr, "invalid input: %v\n", err)
		return apperrors.ExitCode(err)
	}

	listener, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		fmt.Fprintf(stderr, "ui failed: %v\n", apperrors.Wrap(apperrors.InvalidInput, "listen on UI address", err))
		return apperrors.ExitInvalidInput
	}
	defer listener.Close()

	server := &http.Server{
		Handler:           newUIHandler(scanWithConfig),
		ReadHeaderTimeout: 5 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Serve(listener)
	}()

	url := "http://" + listener.Addr().String() + "/"
	fmt.Fprintf(stdout, "VelesMist UI listening on %s\n", url)
	if cfg.Open {
		if err := openBrowser(url); err != nil {
			fmt.Fprintf(stderr, "browser open failed: %v\n", err)
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stderr, "ui shutdown failed: %v\n", err)
			return apperrors.ExitInternal
		}
		return apperrors.ExitSuccess
	case err := <-errCh:
		if err == nil || err == http.ErrServerClosed {
			return apperrors.ExitSuccess
		}
		fmt.Fprintf(stderr, "ui failed: %v\n", err)
		return apperrors.ExitInternal
	}
}

func parseUI(args []string) (uiConfig, error) {
	cfg := uiConfig{
		Addr: defaultUIAddr,
		Open: true,
	}
	fs := flag.NewFlagSet("ui", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&cfg.Addr, "addr", cfg.Addr, "local UI listen address")
	fs.BoolVar(&cfg.Open, "open", cfg.Open, "open the UI in the default browser")
	if err := fs.Parse(args); err != nil {
		return uiConfig{}, apperrors.Wrap(apperrors.InvalidInput, "parse UI flags", err)
	}
	if fs.NArg() != 0 {
		return uiConfig{}, apperrors.New(apperrors.InvalidInput, "ui does not accept positional arguments")
	}
	host, _, err := net.SplitHostPort(cfg.Addr)
	if err != nil {
		return uiConfig{}, apperrors.Wrap(apperrors.InvalidInput, "ui address must be host:port", err)
	}
	if host != "127.0.0.1" && host != "localhost" && host != "::1" {
		return uiConfig{}, apperrors.New(apperrors.InvalidInput, "ui address must bind to localhost")
	}
	return cfg, nil
}

func newUIHandler(scan uiScanFunc) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; script-src 'unsafe-inline'; connect-src 'self'; base-uri 'none'; form-action 'none'")
		_, _ = io.WriteString(w, uiHTML)
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = io.WriteString(w, "ok\n")
	})
	mux.HandleFunc("/api/scan", handleUIScan(scan))
	return mux
}

func handleUIScan(scan uiScanFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeUIError(w, apperrors.New(apperrors.InvalidInput, "method not allowed"))
			return
		}
		defer r.Body.Close()
		var req uiScanRequest
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			writeUIError(w, apperrors.Wrap(apperrors.InvalidInput, "decode scan request", err))
			return
		}
		cfg, err := scanConfigFromUI(req)
		if err != nil {
			writeUIError(w, err)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), cfg.Timeout)
		defer cancel()
		result, err := scan(ctx, cfg)
		if err != nil {
			writeUIError(w, err)
			return
		}
		writeUIJSON(w, http.StatusOK, uiScanResponse{OK: true, Report: result})
	}
}

func scanConfigFromUI(req uiScanRequest) (config.ScanConfig, error) {
	args := []string{"--format", config.FormatJSON}
	game := strings.TrimSpace(req.Game)
	if game == "" {
		game = config.GameDota2
	}
	args = append(args, "--game", game)

	fixture := strings.TrimSpace(req.Fixture)
	if fixture != "" {
		args = append(args, "--fixture", fixture)
	} else {
		args = append(args, "--steam-id", strings.TrimSpace(req.SteamID))
	}
	if priceCache := strings.TrimSpace(req.PriceCache); priceCache != "" {
		args = append(args, "--price-cache", priceCache)
	}
	if minPrice := strings.TrimSpace(req.MinPrice); minPrice != "" {
		args = append(args, "--min-price", minPrice)
	}
	if currency := strings.TrimSpace(req.Currency); currency != "" {
		args = append(args, "--currency", currency)
	}
	if req.NoCache {
		args = append(args, "--no-cache")
	}
	return config.ParseScan(args)
}

func writeUIError(w http.ResponseWriter, err error) {
	kind := apperrors.KindOf(err)
	status := http.StatusInternalServerError
	switch kind {
	case apperrors.InvalidInput:
		status = http.StatusBadRequest
	case apperrors.RateLimited:
		status = http.StatusTooManyRequests
	case apperrors.NetworkTimeout:
		status = http.StatusGatewayTimeout
	case apperrors.Upstream:
		status = http.StatusBadGateway
	}
	writeUIJSON(w, status, uiErrorResponse{
		OK:       false,
		Kind:     string(kind),
		ExitCode: apperrors.ExitCode(err),
		Message:  safeUIErrorMessage(err),
	})
}

func safeUIErrorMessage(err error) string {
	switch apperrors.KindOf(err) {
	case apperrors.InvalidInput:
		return err.Error()
	case apperrors.RateLimited:
		return "Steam rate limited the request. Wait and try again."
	case apperrors.NetworkTimeout:
		return "Steam request timed out. Check network access and try again."
	case apperrors.Upstream:
		return "Steam inventory is unavailable, private, rate limited, or returned an upstream error."
	default:
		return "Internal error while running the scan."
	}
}

func writeUIJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

const uiHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>VelesMist</title>
  <style>
    :root { color-scheme: light dark; --bg:#f6f7f8; --panel:#ffffff; --text:#15181c; --muted:#5f6875; --line:#d8dde4; --accent:#116b5f; --accent-2:#0d534a; --danger:#b42318; --code:#eef1f4; }
    @media (prefers-color-scheme: dark) { :root { --bg:#151719; --panel:#202327; --text:#edf0f2; --muted:#a7b0ba; --line:#343a42; --accent:#53b8a9; --accent-2:#78ccbf; --danger:#ff9b8f; --code:#16191c; } }
    * { box-sizing: border-box; }
    body { margin:0; font:14px/1.45 system-ui, -apple-system, Segoe UI, sans-serif; background:var(--bg); color:var(--text); }
    main { max-width:1180px; margin:0 auto; padding:24px; }
    header { display:flex; justify-content:space-between; gap:16px; align-items:flex-start; margin-bottom:18px; }
    h1 { margin:0; font-size:24px; font-weight:700; letter-spacing:0; }
    .sub { margin-top:4px; color:var(--muted); }
    .grid { display:grid; grid-template-columns: 360px 1fr; gap:16px; align-items:start; }
    section, aside { background:var(--panel); border:1px solid var(--line); border-radius:8px; padding:16px; }
    label { display:block; font-weight:650; margin:12px 0 6px; }
    input, select { width:100%; border:1px solid var(--line); border-radius:6px; background:transparent; color:var(--text); padding:9px 10px; font:inherit; }
    input[type=checkbox] { width:auto; margin-right:8px; }
    .row { display:grid; grid-template-columns: 1fr 1fr; gap:10px; }
    .check { display:flex; align-items:center; color:var(--muted); margin-top:12px; }
    button { border:0; border-radius:6px; background:var(--accent); color:white; padding:10px 14px; font:inherit; font-weight:700; cursor:pointer; width:100%; margin-top:16px; }
    button:hover { background:var(--accent-2); }
    button:disabled { opacity:.6; cursor:wait; }
    .status { min-height:20px; margin-top:12px; color:var(--muted); }
    .error { color:var(--danger); }
    .summary { display:grid; grid-template-columns: repeat(4, minmax(110px, 1fr)); gap:10px; margin-bottom:14px; }
    .metric { border:1px solid var(--line); border-radius:6px; padding:10px; }
    .metric b { display:block; font-size:18px; }
    .tabs { display:flex; gap:8px; margin:8px 0 12px; }
    .tab { width:auto; margin:0; padding:7px 10px; background:transparent; color:var(--text); border:1px solid var(--line); }
    .tab.active { background:var(--accent); color:white; border-color:var(--accent); }
    .table-wrap { overflow:auto; border:1px solid var(--line); border-radius:6px; }
    table { width:100%; border-collapse:collapse; min-width:820px; }
    th, td { text-align:left; border-bottom:1px solid var(--line); padding:9px 10px; vertical-align:top; }
    th { background:rgba(120,120,120,.08); font-size:12px; text-transform:uppercase; color:var(--muted); letter-spacing:.02em; }
    tr:last-child td { border-bottom:0; }
    code, pre { background:var(--code); border-radius:6px; }
    pre { overflow:auto; padding:12px; margin:0; max-height:520px; }
    .pill { display:inline-block; padding:2px 7px; border-radius:999px; border:1px solid var(--line); font-size:12px; }
    .empty { color:var(--muted); padding:40px 0; text-align:center; border:1px dashed var(--line); border-radius:6px; }
    .advanced { margin-top:14px; padding-top:10px; border-top:1px solid var(--line); }
    .hint { color:var(--muted); font-size:12px; margin-top:6px; }
    @media (max-width: 860px) { main { padding:14px; } header, .grid { display:block; } aside { margin-bottom:14px; } .summary { grid-template-columns: repeat(2, minmax(110px, 1fr)); } }
  </style>
</head>
<body>
<main>
  <header>
    <div>
      <h1>VelesMist</h1>
      <div class="sub">Read-only Steam/Dota 2 inventory scan. No login, cookies, listings, or selling.</div>
    </div>
  </header>
  <div class="grid">
    <aside>
      <form id="scan-form">
        <label for="steam-id">Steam64 ID</label>
        <input id="steam-id" name="steam_id" inputmode="numeric" autocomplete="off" placeholder="76561197987179126">
        <div class="hint">Inventory must be public/readable. Private inventories return an upstream error.</div>
        <div class="row">
          <div>
            <label for="game">Game</label>
            <select id="game" name="game"><option value="dota2">Dota 2</option></select>
          </div>
          <div>
            <label for="min-price">Min receive</label>
            <input id="min-price" name="min_price" value="0" placeholder="5.00">
          </div>
        </div>
        <div class="advanced">
          <label for="price-cache">Price cache path</label>
          <input id="price-cache" name="price_cache" placeholder="C:\path\prices.json">
          <label for="fixture">Offline fixture path</label>
          <input id="fixture" name="fixture" placeholder="internal/inventory/testdata/dota_inventory.json">
          <label class="check"><input id="no-cache" name="no_cache" type="checkbox"> Skip inventory cache</label>
        </div>
        <button id="scan-button" type="submit">Scan inventory</button>
        <div id="status" class="status"></div>
      </form>
    </aside>
    <section>
      <div id="output"><div class="empty">Run a scan to see recommendations.</div></div>
    </section>
  </div>
</main>
<script>
const form = document.getElementById('scan-form');
const statusBox = document.getElementById('status');
const output = document.getElementById('output');
const button = document.getElementById('scan-button');
let lastReport = null;

form.addEventListener('submit', async (event) => {
  event.preventDefault();
  statusBox.className = 'status';
  statusBox.textContent = 'Scanning...';
  button.disabled = true;
  const data = Object.fromEntries(new FormData(form).entries());
  data.no_cache = document.getElementById('no-cache').checked;
  try {
    const res = await fetch('/api/scan', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(data)
    });
    const body = await res.json();
    if (!res.ok || !body.ok) {
      throw new Error(body.message || 'scan failed');
    }
    lastReport = body.report;
    statusBox.textContent = 'Scan complete.';
    renderReport(lastReport, 'table');
  } catch (err) {
    statusBox.className = 'status error';
    statusBox.textContent = err.message;
  } finally {
    button.disabled = false;
  }
});

function money(cents, currency) {
  const sign = cents < 0 ? '-' : '';
  const value = Math.abs(cents || 0);
  return currency + ' ' + sign + Math.floor(value / 100) + '.' + String(value % 100).padStart(2, '0');
}

function renderReport(report, mode) {
  const s = report.summary || {};
  const metrics =
    '<div class="summary">' +
      '<div class="metric"><span>Sell</span><b>' + (s.candidate_items || 0) + '</b></div>' +
      '<div class="metric"><span>Priced</span><b>' + (s.priced_items || 0) + '</b></div>' +
      '<div class="metric"><span>Missing price</span><b>' + (s.missing_price_items || 0) + '</b></div>' +
      '<div class="metric"><span>Receive</span><b>' + money(s.estimated_total_receive_cents || 0, report.currency) + '</b></div>' +
    '</div>';
  const tabs =
    '<div class="tabs">' +
      '<button class="tab ' + (mode === 'table' ? 'active' : '') + '" onclick="renderReport(lastReport, \'table\')">Table</button>' +
      '<button class="tab ' + (mode === 'json' ? 'active' : '') + '" onclick="renderReport(lastReport, \'json\')">JSON</button>' +
    '</div>';
  if (mode === 'json') {
    output.innerHTML = metrics + tabs + '<pre>' + escapeHTML(JSON.stringify(report, null, 2)) + '</pre>';
    return;
  }
  const rows = (report.items || []).map(item =>
    '<tr>' +
      '<td>' + escapeHTML(item.market_hash_name || item.name || '') + '</td>' +
      '<td>' + (item.count || 0) + '</td>' +
      '<td>' + (item.price_status === 'priced' ? money(item.buyer_price_cents, report.currency) : '-') + '</td>' +
      '<td>' + (item.price_status === 'priced' ? money(item.seller_receive_cents, report.currency) : '-') + '</td>' +
      '<td>' + (item.price_status === 'priced' ? money(item.total_receive_cents, report.currency) : '-') + '</td>' +
      '<td><span class="pill">' + escapeHTML(item.recommendation || '') + '</span></td>' +
      '<td>' + escapeHTML((item.reason_codes || []).join(', ')) + '</td>' +
      '<td><a href="' + escapeAttr(item.market_url || '#') + '" target="_blank" rel="noreferrer">market</a></td>' +
    '</tr>').join('');
  output.innerHTML = metrics + tabs +
    '<div class="table-wrap"><table>' +
      '<thead><tr><th>Item</th><th>Count</th><th>Gross</th><th>Receive each</th><th>Total</th><th>Status</th><th>Reasons</th><th>Link</th></tr></thead>' +
      '<tbody>' + (rows || '<tr><td colspan="8">No marketable items in report.</td></tr>') + '</tbody>' +
    '</table></div>';
}

function escapeHTML(value) {
  return String(value).replace(/[&<>"']/g, ch => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[ch]));
}

function escapeAttr(value) {
  return escapeHTML(value).replace(/\x60/g, '&#96;');
}
</script>
</body>
</html>`
