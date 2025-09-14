package server

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// keep your existing types
type Config struct {
	DbPath        string
	TemplatesPath string
	StaticPath    string
	QrPath        string
}

var (
	// env var names
	envDbPath        = "BILLDB_DB_PATH"
	envTemplatesPath = "BILLDB_TEMPLATE_PATH"
	envStaticPath    = "BILLDB_STATIC_PATH"
	envQrPath        = "BILLDB_QR_TMP_PATH"
)

// LoadConfig tries CLI flags first, then env vars, then a config file (if provided via CLI).
// It enforces that a single method must supply all required fields; partials are discarded
// and the next method is attempted. If after all methods required fields are missing,
// it prints which are missing to stdout and returns an error.
func LoadConfig() (*Config, error) {
	// 1) parse CLI flags (but do not use partials — we'll validate)
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// CLI flags corresponding to config fields
	cliDbPath := fs.String("db-path", "", "path to DB (BILLDB_DB_PATH)")
	cliTemplatesPath := fs.String("templates-path", "", "path to templates (BILLDB_TEMPLATE_PATH)")
	cliStaticPath := fs.String("static-path", "", "path to static files (BILLDB_STATIC_PATH)")
	cliQrPath := fs.String("qr-path", "", "path to qr tmp (BILLDB_QR_TMP_PATH)")

	// config-file flag: path to KEY=VALUE file
	cliConfigFile := fs.String("config-file", "", "path to config file with KEY=VALUE lines matching env var names")

	// parse the flags; ignore errors to allow program to continue returning useful errors later
	_ = fs.Parse(os.Args[1:])

	// Helper to collect missing keys
	missing := func(c *Config) []string {
		var miss []string
		if c.DbPath == "" {
			miss = append(miss, envDbPath)
		}
		if c.TemplatesPath == "" {
			miss = append(miss, envTemplatesPath)
		}
		if c.StaticPath == "" {
			miss = append(miss, envStaticPath)
		}
		if c.QrPath == "" {
			miss = append(miss, envQrPath)
		}
		return miss
	}

	// Try 1: CLI flags (must be complete)
	cliCfg := &Config{
		DbPath:        strings.TrimSpace(*cliDbPath),
		TemplatesPath: strings.TrimSpace(*cliTemplatesPath),
		StaticPath:    strings.TrimSpace(*cliStaticPath),
		QrPath:        strings.TrimSpace(*cliQrPath),
	}

	if len(missing(cliCfg)) == 0 {
		// full config provided by CLI flags
		return cliCfg, nil
	}

	// If CLI flags were partially provided, discard them (clean config) before trying env.
	// (Per requirement: "if cli flags are partial, before reading env vars clean config flags if any present after cli")
	// So we create a fresh config when reading env.
	envCfg := &Config{}

	// Try 2: environment variables (must be complete)
	if v, ok := os.LookupEnv(envDbPath); ok {
		envCfg.DbPath = strings.TrimSpace(v)
	}
	if v, ok := os.LookupEnv(envTemplatesPath); ok {
		envCfg.TemplatesPath = strings.TrimSpace(v)
	}
	if v, ok := os.LookupEnv(envStaticPath); ok {
		envCfg.StaticPath = strings.TrimSpace(v)
	}
	if v, ok := os.LookupEnv(envQrPath); ok {
		envCfg.QrPath = strings.TrimSpace(v)
	}

	if len(missing(envCfg)) == 0 {
		return envCfg, nil
	}

	// If env is partial, discard and try config-file if provided via CLI.
	// Per requirement: "firstly read cli flags, then if config is empty or partial read env vars, at the end if config is empty or partial send to stdout error"
	// We already tried CLI and env. Now try config-file only if CLI flag --config-file was provided.
	if *cliConfigFile != "" {
		fileCfg := &Config{}
		if err := readConfigFile(*cliConfigFile, fileCfg); err != nil {
			// If reading the file fails, report which sources we tried and error
			fmt.Fprintf(os.Stdout, "failed to read config file %q: %v\n", *cliConfigFile, err)
			// fallthrough to final error reporting below
		} else {
			if len(missing(fileCfg)) == 0 {
				return fileCfg, nil
			}
		}
	}

	// Nothing produced a full config. Print which fields are missing from each attempted method (preference order).
	// Determine the final best "partial" to report missing keys from: prefer CLI if it had any non-empty, else env, else config-file.
	var reportCfg *Config
	var method string
	if anySet(cliCfg) {
		reportCfg = cliCfg
		method = "CLI flags"
	} else if anySet(envCfg) {
		reportCfg = envCfg
		method = "environment variables"
	} else if *cliConfigFile != "" {
		// attempt to read file once more into a fresh struct for reporting; ignore read error
		fc := &Config{}
		_ = readConfigFile(*cliConfigFile, fc)
		reportCfg = fc
		method = fmt.Sprintf("config file (%s)", *cliConfigFile)
	} else {
		reportCfg = &Config{}
		method = "no configuration provided"
	}

	miss := missing(reportCfg)
	// Print missing keys to stdout as requested
	if len(miss) > 0 {
		fmt.Fprintln(os.Stdout, "Missing configuration fields (tried", method+"):")
		for _, k := range miss {
			fmt.Fprintln(os.Stdout, "- "+k)
		}
	} else {
		// unlikely: if none missing then return success (defensive)
		return reportCfg, nil
	}

	return nil, errors.New("incomplete configuration")
}

// anySet returns true if any field in cfg is non-empty
func anySet(cfg *Config) bool {
	if cfg == nil {
		return false
	}
	return cfg.DbPath != "" || cfg.TemplatesPath != "" || cfg.StaticPath != "" || cfg.QrPath != ""
}

// readConfigFile reads KEY=VALUE lines from path and populates cfg.
// Recognizes the same keys as env var names (BILLDB_DB_PATH, BILLDB_TEMPLATE_PATH, etc.).
// Lines starting with # are treated as comments. Empty values are permitted but will be set as empty strings.
func readConfigFile(path string, cfg *Config) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Split only on first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line %d: %q", lineNum, line)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case envDbPath:
			cfg.DbPath = val
		case envTemplatesPath:
			cfg.TemplatesPath = val
		case envStaticPath:
			cfg.StaticPath = val
		case envQrPath:
			cfg.QrPath = val
		default:
			// ignore unknown keys
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
