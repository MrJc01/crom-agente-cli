package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:     "open",
	Aliases: []string{"app"},
	Short:   "Abre o aplicativo de desktop do crom-agente (Tauri app)",
	Long:    `Busca e inicializa o aplicativo de desktop do crom-agente se ele estiver compilado ou disponível no sistema.`,
	Example: `  $ crom-agente-cli open
  $ crom-agente-cli app`,
	RunE: func(cmd *cobra.Command, args []string) error {
		binName := "crom-agente-app"
		if runtime.GOOS == "windows" {
			binName = "crom-agente-app.exe"
		}

		var candidates []string

		// 1. Procurar no PATH
		if appPath, err := exec.LookPath(binName); err == nil {
			candidates = append(candidates, appPath)
		}

		// 2. Caminhos relativos a partir do CWD
		if cwd, err := os.Getwd(); err == nil {
			candidates = append(candidates,
				filepath.Join(cwd, "crom-agente-app", "src-tauri", "target", "release", binName),
				filepath.Join(cwd, "crom-agente-app", "src-tauri", "target", "debug", binName),
				filepath.Join(cwd, "..", "crom-agente-app", "src-tauri", "target", "release", binName),
				filepath.Join(cwd, "..", "crom-agente-app", "src-tauri", "target", "debug", binName),
			)
		}

		// 3. Caminhos relativos a partir do executável do CLI
		if cliExe, err := os.Executable(); err == nil {
			cliDir := filepath.Dir(cliExe)
			candidates = append(candidates,
				filepath.Join(cliDir, "crom-agente-app", "src-tauri", "target", "release", binName),
				filepath.Join(cliDir, "crom-agente-app", "src-tauri", "target", "debug", binName),
				filepath.Join(cliDir, "..", "crom-agente-app", "src-tauri", "target", "release", binName),
				filepath.Join(cliDir, "..", "crom-agente-app", "src-tauri", "target", "debug", binName),
			)
		}

		// 4. Caminho padrão na pasta do usuário (Documentos/GitHub)
		homeDir, err := os.UserHomeDir()
		if err == nil {
			candidates = append(candidates,
				filepath.Join(homeDir, "Documentos", "GitHub", "crom-agente-app", "src-tauri", "target", "release", binName),
				filepath.Join(homeDir, "Documentos", "GitHub", "crom-agente-app", "src-tauri", "target", "debug", binName),
			)
		}

		// Filtrar duplicados e caminhos inexistentes
		var foundPath string
		for _, path := range candidates {
			if path == "" {
				continue
			}
			if info, err := os.Stat(path); err == nil && !info.IsDir() {
				foundPath = path
				break
			}
		}

		// Se não encontrou o binário, tentar fallback de execução dev caso esteja no workspace
		if foundPath == "" {
			// Procurar se a pasta de desenvolvimento do app com package.json existe
			var packageDirs []string
			if cwd, err := os.Getwd(); err == nil {
				packageDirs = append(packageDirs,
					filepath.Join(cwd, "crom-agente-app"),
					filepath.Join(cwd, "..", "crom-agente-app"),
				)
			}
			if homeDir != "" {
				packageDirs = append(packageDirs,
					filepath.Join(homeDir, "Documentos", "GitHub", "crom-agente-app"),
				)
			}

			for _, dir := range packageDirs {
				pkgJson := filepath.Join(dir, "package.json")
				if info, err := os.Stat(pkgJson); err == nil && !info.IsDir() {
					fmt.Println("Binário compilado não encontrado. Iniciando via 'npm run tauri dev'...")
					
					runCmd := exec.Command("npm", "run", "tauri", "dev")
					runCmd.Dir = dir
					
					err = startCommand(runCmd)
					if err != nil {
						return fmt.Errorf("falha ao iniciar tauri dev: %w", err)
					}
					fmt.Println("🚀 Aplicativo iniciando em modo desenvolvimento no background.")
					return nil
				}
			}

			return fmt.Errorf("aplicativo desktop '%s' não encontrado. Certifique-se de que o projeto foi compilado ou que a pasta 'crom-agente-app' está no diretório correto", binName)
		}

		fmt.Printf("Iniciando aplicativo a partir de: %s\n", foundPath)

		// Executa desassociado usando o helper startCommand
		runCmd := exec.Command(foundPath)
		err = startCommand(runCmd)
		if err != nil {
			return fmt.Errorf("erro ao abrir o aplicativo: %w", err)
		}

		fmt.Println("🚀 Aplicativo aberto com sucesso em segundo plano!")
		return nil
	},
}
