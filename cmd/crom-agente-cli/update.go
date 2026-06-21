package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Atualiza o CLI e o daemon crom-agente para a última versão disponível",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Verificando atualizações...")

		// Identificar OS e Arch
		osName := runtime.GOOS
		arch := runtime.GOARCH

		var agentBinName, cliBinName string
		if osName == "linux" && arch == "amd64" {
			agentBinName = "crom-agente-linux-amd64"
			cliBinName = "crom-agente-cli-linux-amd64"
		} else if osName == "linux" && arch == "arm64" {
			agentBinName = "crom-agente-linux-arm64"
			cliBinName = "crom-agente-cli-linux-arm64"
		} else if osName == "darwin" && arch == "amd64" {
			agentBinName = "crom-agente-darwin-amd64"
			cliBinName = "crom-agente-cli-darwin-amd64"
		} else if osName == "darwin" && arch == "arm64" {
			agentBinName = "crom-agente-darwin-arm64"
			cliBinName = "crom-agente-cli-darwin-arm64"
		} else {
			return fmt.Errorf("arquitetura não suportada pelo auto-update: %s/%s", osName, arch)
		}

		// 1. Atualizar o Agente (se instalado)
		agentPath, err := exec.LookPath("crom-agente")
		if err == nil {
			fmt.Printf("Agente encontrado em: %s\n", agentPath)
			err = updateBinary("MrJc01/crom-agente", agentBinName, agentPath)
			if err != nil {
				if strings.Contains(err.Error(), "permission denied") {
					fmt.Println("\n❌ Permissão negada ao atualizar o crom-agente.")
					fmt.Println("Por favor, execute o comando com sudo: sudo crom-agente-cli update")
					return err
				}
				fmt.Printf("Aviso: Falha ao atualizar agente: %v\n", err)
			} else {
				fmt.Println("✅ crom-agente atualizado com sucesso!")
			}
		} else {
			fmt.Println("ℹ️ crom-agente não encontrado no PATH. Pulando atualização do daemon.")
		}

		// 2. Atualizar o próprio CLI
		cliPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("não foi possível determinar o caminho do CLI: %v", err)
		}
		
		fmt.Printf("CLI encontrado em: %s\n", cliPath)
		err = updateBinary("MrJc01/crom-agente-cli", cliBinName, cliPath)
		if err != nil {
			if strings.Contains(err.Error(), "permission denied") {
				fmt.Println("\n❌ Permissão negada ao atualizar o crom-agente-cli.")
				fmt.Println("Por favor, execute o comando com sudo: sudo crom-agente-cli update")
				return err
			}
			return fmt.Errorf("falha ao atualizar CLI: %v", err)
		}

		fmt.Println("✅ crom-agente-cli atualizado com sucesso!")
		return nil
	},
}

func updateBinary(repo string, assetName string, destPath string) error {
	// Buscar latest release
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo))
	if err != nil {
		return fmt.Errorf("erro ao conectar no github: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github api retornou status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("falha ao decodificar JSON da release: %v", err)
	}

	var downloadUrl string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadUrl = asset.BrowserDownloadURL
			break
		}
	}

	if downloadUrl == "" {
		return fmt.Errorf("binário %s não encontrado na release %s", assetName, release.TagName)
	}

	fmt.Printf("Baixando %s versão %s...\n", assetName, release.TagName)

	// Fazer download para arquivo temporário
	tmpFile := filepath.Join(os.TempDir(), assetName+"_update")
	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	dlResp, err := http.Get(downloadUrl)
	if err != nil {
		out.Close()
		return err
	}
	defer dlResp.Body.Close()

	if dlResp.StatusCode != http.StatusOK {
		out.Close()
		return fmt.Errorf("falha ao baixar binário, HTTP status %d", dlResp.StatusCode)
	}

	_, err = io.Copy(out, dlResp.Body)
	out.Close()
	if err != nil {
		return err
	}

	// Dar permissão de execução ao tmp
	if err := os.Chmod(tmpFile, 0755); err != nil {
		return err
	}

	// Substituir binário
	// No linux/mac é seguro usar rename mesmo com o binário rodando (no caso do CLI)
	// Em alguns sistemas pode ser necessário remover o original antes.
	if runtime.GOOS != "windows" {
		err = os.Rename(tmpFile, destPath)
		if err != nil {
			// Tentar fallback copiando se rename cruzar partições
			if errCopy := copyFile(tmpFile, destPath); errCopy != nil {
				return errCopy // retorna o erro da cópia que normalmente diz permission denied
			}
		}
	} else {
		// Windows: rename para .old antes
		os.Rename(destPath, destPath+".old")
		if err := os.Rename(tmpFile, destPath); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
