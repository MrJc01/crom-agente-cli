package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

// getAppFileName retorna o nome do arquivo do aplicativo de acordo com o OS.
func getAppFileName() string {
	switch runtime.GOOS {
	case "darwin":
		return "crom-agente.dmg"
	case "windows":
		return "crom-agente-app.exe"
	default: // linux
		return "crom-agente.AppImage"
	}
}

// getDevBinName retorna o nome do binário compilado (dev build) do Tauri.
func getDevBinName() string {
	if runtime.GOOS == "windows" {
		return "crom-agente-app.exe"
	}
	return "crom-agente-app"
}

var openCmd = &cobra.Command{
	Use:     "open",
	Aliases: []string{"app"},
	Short:   "Abre o aplicativo de desktop do crom-agente (Tauri app)",
	Long:    `Busca e inicializa o aplicativo de desktop do crom-agente se ele estiver instalado via install.sh ou compilado localmente.`,
	Example: `  $ crom-agente-cli open
  $ crom-agente-cli app`,
	RunE: func(cmd *cobra.Command, args []string) error {
		appFile := getAppFileName()
		devBin := getDevBinName()

		homeDir, _ := os.UserHomeDir()

		var candidates []string

		// ────────────────────────────────────────────────────────
		// 1. Caminhos de instalação do install.sh (prioridade máxima)
		//    O install.sh coloca o AppImage/DMG em:
		//    ~/Desktop/CromIA/ ou ~/CromIA/ (se Desktop não existir)
		// ────────────────────────────────────────────────────────
		if homeDir != "" {
			candidates = append(candidates,
				filepath.Join(homeDir, "Desktop", "CromIA", appFile),
				filepath.Join(homeDir, "Área de Trabalho", "CromIA", appFile),
				filepath.Join(homeDir, "CromIA", appFile),
			)
		}

		// ────────────────────────────────────────────────────────
		// 2. Procurar o binário do app no PATH do sistema
		// ────────────────────────────────────────────────────────
		if appPath, err := exec.LookPath(appFile); err == nil {
			candidates = append(candidates, appPath)
		}
		if appPath, err := exec.LookPath(devBin); err == nil {
			candidates = append(candidates, appPath)
		}

		// ────────────────────────────────────────────────────────
		// 3. Caminhos relativos a partir do executável do CLI
		//    (ex: /usr/local/bin/crom-agente.AppImage)
		// ────────────────────────────────────────────────────────
		if cliExe, err := os.Executable(); err == nil {
			cliDir := filepath.Dir(cliExe)
			candidates = append(candidates,
				filepath.Join(cliDir, appFile),
			)
		}

		// ────────────────────────────────────────────────────────
		// 4. Caminhos de compilação local (dev builds Tauri)
		// ────────────────────────────────────────────────────────
		if cwd, err := os.Getwd(); err == nil {
			candidates = append(candidates,
				filepath.Join(cwd, "crom-agente-app", "src-tauri", "target", "release", devBin),
				filepath.Join(cwd, "crom-agente-app", "src-tauri", "target", "debug", devBin),
				filepath.Join(cwd, "..", "crom-agente-app", "src-tauri", "target", "release", devBin),
				filepath.Join(cwd, "..", "crom-agente-app", "src-tauri", "target", "debug", devBin),
			)
		}

		if homeDir != "" {
			candidates = append(candidates,
				filepath.Join(homeDir, "Documentos", "GitHub", "crom-agente-app", "src-tauri", "target", "release", devBin),
				filepath.Join(homeDir, "Documentos", "GitHub", "crom-agente-app", "src-tauri", "target", "debug", devBin),
			)
		}

		// ────────────────────────────────────────────────────────
		// Verificar cada candidato na ordem de prioridade
		// ────────────────────────────────────────────────────────
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

		// ────────────────────────────────────────────────────────
		// Fallback: npm run tauri dev (ambiente de desenvolvimento)
		// ────────────────────────────────────────────────────────
		if foundPath == "" {
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

			return fmt.Errorf("aplicativo desktop não encontrado.\n\n" +
				"O CLI procurou nos seguintes locais:\n" +
				"  • ~/Desktop/CromIA/" + appFile + "\n" +
				"  • ~/CromIA/" + appFile + "\n" +
				"  • PATH do sistema\n" +
				"  • Compilações locais (src-tauri/target/)\n\n" +
				"Para instalar, execute:\n" +
				"  curl -sSL https://cloud.ia.crom.run/install.sh | bash -s app")
		}

		fmt.Printf("✅ Aplicativo encontrado: %s\n", foundPath)

		// No macOS, usar "open" para .dmg
		if runtime.GOOS == "darwin" && filepath.Ext(foundPath) == ".dmg" {
			fmt.Println("Montando o instalador DMG...")
			runCmd := exec.Command("open", foundPath)
			err := startCommand(runCmd)
			if err != nil {
				return fmt.Errorf("erro ao abrir o DMG: %w", err)
			}
			fmt.Println("🚀 Instalador DMG aberto com sucesso!")
			return nil
		}

		// Executa desassociado usando o helper startCommand
		runCmd := exec.Command(foundPath)
		err := startCommand(runCmd)
		if err != nil {
			return fmt.Errorf("erro ao abrir o aplicativo: %w", err)
		}

		fmt.Println("🚀 Aplicativo aberto com sucesso em segundo plano!")
		return nil
	},
}

