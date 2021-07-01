package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"globetrot/database"
	"hash"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"time"
)

type Runner struct {
	Config      Config
	Database    database.Database
	Hasher      hash.Hash
	Logger      *Logger
	FileManager *FileManager
}

func (r *Runner) Init(config Config) {
	r.Config = config
	r.Database = database.NewDatabase(config.Type)
	r.Hasher = sha256.New()
	r.FileManager = &FileManager{filePath: config.FilePath}
	r.Logger = &Logger{debug: true}
}

func (r *Runner) Migrate() {

	start := time.Now()

	r.Logger.PrintConfig(r.Config)

	// Ensure proper directory structure
	err := r.FileManager.EnsureDirectory()
	if err != nil {
		r.Logger.Fatal(err.Error())
	}

	// Connect to the database
	r.Database.Connect(r.Config.Username, r.Config.Password, r.Config.Host, r.Config.Port, r.Config.Database)

	// Generate migration metadata table
	r.Database.CreateMigrationsTable()
	r.Logger.Details("Creating migration metadata table if necessary.\n")

	upPath := r.Config.FilePath + "/up"
	p := filepath.FromSlash(upPath)

	files, err := r.FileManager.GetUpScripts()
	if err != nil {
		r.Logger.Fatal(err.Error())
	}

	for _, f := range files {
		r.RunUpScript(p, f)
	}

	r.Database.Close()

	t := time.Now()
	elapsed := t.Sub(start)

	r.Logger.Success(fmt.Sprintf("\n-----------------------------\nMigration complete. %v elapsed\n-----------------------------", elapsed))
}

func (r Runner) RunUpScript(upPath string, file fs.FileInfo) {
	script_name := file.Name()

	b, err := ioutil.ReadFile(upPath + "/" + file.Name())
	if err != nil {
		r.Logger.Fatal(err.Error())
	}

	r.Hasher.Write(b)
	sha := base64.URLEncoding.EncodeToString((r.Hasher.Sum(nil)))
	sql := string(b)
	script_row := r.Database.GetScriptRun(script_name)
	if script_row != nil && script_row.Hash != sha {
		r.Logger.Fatal(fmt.Sprintf("Changing of one-time script '%s' not allowed after application.\n", script_name))
	}
	if script_row == nil {
		r.Database.ApplyScript(sql, script_name, sha)
		r.Logger.Success(fmt.Sprintf("APPLIED SCRIPT: %s", script_name))
	} else {
		r.Logger.Details(fmt.Sprintf("SKIPPING SCRIPT: %s", script_name))
	}
}