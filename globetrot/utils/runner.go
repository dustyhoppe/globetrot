package utils

import (
	"fmt"
	"globetrot/database"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Runner struct {
	Config      Config
	Database    database.Database
	Hasher      *HashGenerator
	Logger      *Logger
	FileManager *FileManager
}

type RunScriptFunc func(path string)
type GetScriptFilesFunc func() ([]string, error)

func (r *Runner) Init(config Config) {
	r.Config = config
	r.Database = database.NewDatabase(config.Type, config.Username, config.Password, config.Host, config.Port, config.Database)

	r.FileManager = &FileManager{filePath: config.FilePath}
	r.Logger = &Logger{debug: true}

	r.Hasher = new(HashGenerator)
	r.Hasher.Init(true)
}

func (r *Runner) Initialize(directory string) {
	r.Logger.OutputHeader()

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.Mkdir(directory, os.ModePerm)
	}

	upPath := filepath.Join(directory, "up")
	if _, err := os.Stat(upPath); os.IsNotExist(err) {
		os.Mkdir(upPath, os.ModePerm)
	}

	procPath := filepath.Join(directory, "procs")
	if _, err := os.Stat(procPath); os.IsNotExist(err) {
		os.Mkdir(procPath, os.ModePerm)
	}

	r.Logger.Success("Initialization complete")
}

func (r *Runner) Migrate() {
	r.Logger.OutputHeader()

	start := time.Now()

	// Ensure proper directory structure
	err := r.FileManager.EnsureDirectory()
	if err != nil {
		r.Logger.Fatal(err.Error())
	}

	// Generate migration metadata table
	r.Database.CreateMigrationsTable()
	r.Logger.Details("Creating migration metadata table if necessary.\n")

	r.applyScripts(r.RunUpScript, r.FileManager.GetUpScripts)
	r.applyScripts(r.RunProcScript, r.FileManager.GetProcScripts)

	t := time.Now()
	elapsed := t.Sub(start)

	r.Logger.Success(fmt.Sprintf("\n-----------------------------\nMigration complete. %v elapsed\n-----------------------------", elapsed))
}

func (r Runner) RunUpScript(upPath string) {
	_, script_name := filepath.Split(upPath)

	b, err := ioutil.ReadFile(upPath)
	if err != nil {
		r.Logger.Fatal(err.Error())
	}

	sql := string(b)
	sha := r.Hasher.GenerateHash(sql)

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

	sql := string(b)
	sha := r.Hasher.GenerateHash(sql)

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
