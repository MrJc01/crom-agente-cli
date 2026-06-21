package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var diagCmd = &cobra.Command{
	Use:   "diag",
	Short: "Executa um diagnóstico do ambiente local (Item 44)",
	Long: `O comando diag verifica se as dependências essenciais do host estão 
instaladas e disponíveis no PATH, incluindo git, docker, go, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Diagnóstico do Ambiente crom-agente")
		fmt.Println("═══════════════════════════════════════")
		fmt.Printf("S.O: %s\n", runtime.GOOS)
		fmt.Printf("Arquitetura: %s\n", runtime.GOARCH)
		fmt.Println("═══════════════════════════════════════")
		
		deps := []string{"git", "docker", "go", "curl", "tar"}
		allGood := true
		
		for _, dep := range deps {
			path, err := exec.LookPath(dep)
			if err != nil {
				fmt.Printf("❌ %-10s : Não encontrado no PATH\n", dep)
				allGood = false
			} else {
				// Tentar pegar versão
				out, _ := exec.Command(dep, "--version").CombinedOutput()
				version := strings.Split(strings.TrimSpace(string(out)), "\n")[0]
				if version == "" {
					version = path
				}
				// truncar versao se for mt longa
				if len(version) > 40 {
					version = version[:40] + "..."
				}
				fmt.Printf("✅ %-10s : %s\n", dep, version)
			}
		}
		
		fmt.Println("═══════════════════════════════════════")
		if allGood {
			fmt.Println("✨ Todas as dependências essenciais estão instaladas!")
		} else {
			fmt.Println("⚠️  Algumas dependências estão faltando. O agente pode ter capacidades limitadas.")
		}
	},
}
