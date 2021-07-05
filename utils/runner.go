package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"globetrot/database"
	"hash"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

type Runner struct {
	Config      Config
	Database    database.Database
	Hasher      hash.Hash
	Logger      *Logger
	FileManager *FileManager
}

type RunScriptFunc func(path string)
type GetScriptFilesFunc func() ([]string, error)

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

	r.applyScripts(r.RunUpScript, r.FileManager.GetUpScripts)
	r.applyScripts(r.RunProcScript, r.FileManager.GetProcScripts)

	r.Database.Close()

	t := time.Now()
	elapsed := t.Sub(start)

	r.Logger.Success(fmt.Sprintf("\n-----------------------------\nMigration complete. %v elapsed\n-----------------------------", elapsed))
}

func (r Runner) RunUpScript(upPath string) {
	_, script_name := path.Split(upPath)

	b, err := ioutil.ReadFile(upPath)
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
		if !r.Config.DryRun {
			r.Database.ApplyScript(sql, script_name, sha)
		}
		r.Logger.Success(fmt.Sprintf("APPLIED SCRIPT: %s", script_name))
	} else {
		r.Logger.Details(fmt.Sprintf("SKIPPING SCRIPT: %s", script_name))
	}
}

func (r Runner) RunProcScript(procPath string) {
	_, script_name := path.Split(procPath)

	b, err := ioutil.ReadFile(procPath)
	if err != nil {
		r.Logger.Fatal(err.Error())
	}

	r.Hasher.Write(b)
	sha := base64.URLEncoding.EncodeToString((r.Hasher.Sum(nil)))
	sql := string(b)
	script_row := r.Database.GetScriptRun(script_name)

	// proc has changed, re-run the script
	if (script_row != nil && script_row.Hash != sha) || script_row == nil {
		if !r.Config.DryRun {
			r.Database.ApplyScript(sql, script_name, sha)
		}
		r.Logger.Success(fmt.Sprintf("APPLIED PROC: %s", script_name))
	} else {
		r.Logger.Details(fmt.Sprintf("SKIPPING PROC: %s", script_name))
	}
}

func (r Runner) applyScripts(applyFunc RunScriptFunc, getScriptFilesFunc GetScriptFilesFunc) {
	files, err := getScriptFilesFunc()
	if err != nil {
		r.Logger.Fatal(err.Error())
	}

	for _, f := range files {
		if r.shouldApplyScript(f, r.Config.Environment) {
			applyFunc(f)
		}
	}
}

func (r Runner) shouldApplyScript(script_name string, environment string) bool {
	parts := strings.Split(script_name, ".")
	if len(parts) <= 3 {
		return true
	}

	length := len(parts)
	if strings.ToLower(parts[length-2]) != "env" {
		return true
	}

	target_environment := parts[length-3]

	return strings.ToLower(environment) == strings.ToLower(target_environment)
}
