package main

import (
    "database/sql"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "path/filepath"
    "sort"

    _ "github.com/lib/pq"
    "keerja-backend/internal/config"
)

// simple migrator: runs all .up.sql files (or .down.sql) in database/migrations
func main() {
    direction := flag.String("dir", "up", "migration direction: up or down")
    migrationsPath := flag.String("path", "database/migrations", "path to migrations directory")
    flag.Parse()

    cfg := config.LoadConfig()

    // open DB using standard library (pq)
    db, err := sql.Open("postgres", cfg.GetDSN())
    if err != nil {
        log.Fatalf("failed to open db: %v", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatalf("failed to ping db: %v", err)
    }

    files, err := ioutil.ReadDir(*migrationsPath)
    if err != nil {
        log.Fatalf("failed to read migrations dir: %v", err)
    }

    var targets []string
    for _, f := range files {
        if f.IsDir() {
            continue
        }
        name := f.Name()
        // Accept files that end with .up.sql or .down.sql depending on direction
        if *direction == "up" && len(name) > 7 && name[len(name)-7:] == ".up.sql" {
            targets = append(targets, filepath.Join(*migrationsPath, name))
        }
        if *direction == "down" && len(name) > 9 && name[len(name)-9:] == ".down.sql" {
            targets = append(targets, filepath.Join(*migrationsPath, name))
        }
    }

    if len(targets) == 0 {
        log.Println("no migration files found for direction", *direction)
        return
    }

    // sort files for up (ascending) and reverse for down
    sort.Strings(targets)
    if *direction == "down" {
        // reverse order
        for i := 0; i < len(targets)/2; i++ {
            j := len(targets) - 1 - i
            targets[i], targets[j] = targets[j], targets[i]
        }
    }

    for _, file := range targets {
        log.Printf("applying %s\n", file)
        content, err := ioutil.ReadFile(file)
        if err != nil {
            log.Fatalf("failed to read file %s: %v", file, err)
        }

        // Execute individual statements so we can ignore "already exists" errors for idempotency
        stmts := splitSQLStatements(string(content))
        for _, s := range stmts {
            if s == "" {
                continue
            }
            if _, err := db.Exec(s); err != nil {
                // If error indicates object already exists, log and continue
                msg := err.Error()
                if containsAlreadyExists(msg) {
                    log.Printf("warning: statement skipped (already exists): %v", msg)
                    continue
                }
                log.Fatalf("failed to execute statement in %s: %v\nstatement: %s", file, err, s)
            }
        }
    }

    fmt.Printf("migrations %s completed (%d files)\n", *direction, len(targets))
}

// splitSQLStatements splits SQL by semicolon but respects single quotes, double quotes and dollar-quoted strings
func splitSQLStatements(sqlText string) []string {
    var stmts []string
    var cur []rune
    inSingle := false
    inDouble := false
    var dollarTag string

    runes := []rune(sqlText)
    for i := 0; i < len(runes); i++ {
        r := runes[i]
        // detect start of dollar-quote
        if r == '$' && !inSingle && !inDouble {
            // try to read tag
            j := i + 1
            for j < len(runes) && ((runes[j] >= 'a' && runes[j] <= 'z') || (runes[j] >= 'A' && runes[j] <= 'Z') || (runes[j] >= '0' && runes[j] <= '9') || runes[j] == '_') {
                j++
            }
            if j < len(runes) && runes[j] == '$' {
                // we have a dollar-tag start
                tag := string(runes[i : j+1]) // includes both $
                if dollarTag == "" {
                    dollarTag = tag
                    cur = append(cur, runes[i])
                    continue
                } else if dollarTag == tag {
                    // closing tag
                    dollarTag = ""
                    cur = append(cur, runes[i])
                    continue
                }
            }
        }

        // if we're inside a dollar-quote, just append
        if dollarTag != "" {
            cur = append(cur, r)
            continue
        }

        if r == '\'' && !inDouble {
            inSingle = !inSingle
            cur = append(cur, r)
            continue
        }
        if r == '"' && !inSingle {
            inDouble = !inDouble
            cur = append(cur, r)
            continue
        }

        if r == ';' && !inSingle && !inDouble && dollarTag == "" {
            // statement terminator
            s := string(cur)
            if trim(s) != "" {
                stmts = append(stmts, s)
            }
            cur = []rune{}
            continue
        }

        cur = append(cur, r)
    }

    // trailing
    if trim(string(cur)) != "" {
        stmts = append(stmts, string(cur))
    }
    return stmts
}

func trim(s string) string {
    // simple whitespace trim
    i := 0
    j := len(s) - 1
    for i <= j {
        if s[i] == ' ' || s[i] == '\n' || s[i] == '\r' || s[i] == '\t' {
            i++
            continue
        }
        break
    }
    for j >= i {
        if s[j] == ' ' || s[j] == '\n' || s[j] == '\r' || s[j] == '\t' {
            j--
            continue
        }
        break
    }
    if i > j {
        return ""
    }
    return s[i : j+1]
}

func containsAlreadyExists(msg string) bool {
    if msg == "" {
        return false
    }
    if indexOf(msg, "already exists") >= 0 {
        return true
    }
    if indexOf(msg, "duplicate key value") >= 0 {
        return true
    }
    return false
}

func indexOf(s, sub string) int {
    for i := 0; i+len(sub) <= len(s); i++ {
        if s[i:i+len(sub)] == sub {
            return i
        }
    }
    return -1
}
