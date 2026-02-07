package evaluator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xrspace/zerglang/runtime/lexer"
	"github.com/xrspace/zerglang/runtime/parser"
)

// ModuleLoader handles loading and caching of modules.
type ModuleLoader struct {
	cache       map[string]*UserModule // path -> module
	searchPaths []string               // directories to search for modules
	currentDir  string                 // current working directory for relative imports
}

// DefaultLoader is the global module loader instance.
var DefaultLoader *ModuleLoader

func init() {
	DefaultLoader = NewModuleLoader()
}

// NewModuleLoader creates a new module loader.
func NewModuleLoader() *ModuleLoader {
	cwd, _ := os.Getwd()
	return &ModuleLoader{
		cache:       make(map[string]*UserModule),
		searchPaths: []string{cwd},
		currentDir:  cwd,
	}
}

// SetCurrentDir sets the current directory for relative imports.
func (ml *ModuleLoader) SetCurrentDir(dir string) {
	ml.currentDir = dir
}

// AddSearchPath adds a directory to the module search paths.
func (ml *ModuleLoader) AddSearchPath(path string) {
	ml.searchPaths = append(ml.searchPaths, path)
}

// LoadModule loads a module by its import path.
func (ml *ModuleLoader) LoadModule(importPath string) (*UserModule, *Error) {
	// Check cache first
	if mod, ok := ml.cache[importPath]; ok {
		return mod, nil
	}

	// Find the module file
	filePath, err := ml.resolveModulePath(importPath)
	if err != nil {
		return nil, &Error{Message: err.Error()}
	}

	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, &Error{Message: "failed to read module: " + err.Error()}
	}

	// Parse the file
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return nil, &Error{Message: "parse error in module: " + strings.Join(p.Errors(), "; ")}
	}

	// Create a new environment for the module with builtins
	moduleEnv := NewEnvironmentWithBuiltins()

	// Save and set the current directory for nested imports
	oldDir := ml.currentDir
	ml.currentDir = filepath.Dir(filePath)

	// Evaluate the program in the module's environment
	result := Eval(program, moduleEnv)

	// Restore the current directory
	ml.currentDir = oldDir

	if isError(result) {
		return nil, result.(*Error)
	}

	// Extract the module name from the import path
	moduleName := filepath.Base(importPath)
	if ext := filepath.Ext(moduleName); ext != "" {
		moduleName = moduleName[:len(moduleName)-len(ext)]
	}

	// Create the module object
	mod := &UserModule{
		Name: moduleName,
		Env:  moduleEnv,
	}

	// Cache it
	ml.cache[importPath] = mod

	return mod, nil
}

// resolveModulePath finds the actual file path for an import path.
func (ml *ModuleLoader) resolveModulePath(importPath string) (string, error) {
	// Try different extensions
	extensions := []string{"", ".zg", ".zerg"}

	// First, try relative to current directory
	for _, ext := range extensions {
		path := filepath.Join(ml.currentDir, importPath+ext)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Then try each search path
	for _, searchPath := range ml.searchPaths {
		for _, ext := range extensions {
			path := filepath.Join(searchPath, importPath+ext)
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}
	}

	return "", fmt.Errorf("module not found: %s", importPath)
}

// evalImportStatement evaluates an import statement.
func evalImportStatement(node *parser.ImportStatement, env *Environment) Object {
	// Load the module
	mod, err := DefaultLoader.LoadModule(node.Path)
	if err != nil {
		return err
	}

	// Determine the name to bind the module to
	name := node.Alias
	if name == "" {
		// Use the last part of the path as the module name
		name = mod.Name
	}

	// Bind the module to the environment
	env.Declare(name, mod, false)

	return NULL
}
