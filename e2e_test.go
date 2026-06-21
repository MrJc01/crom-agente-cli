package cli_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestE2E_CliDiag(t *testing.T) {
	// Item 45: Testes E2E rodando comandos inteiros da CLI simulados
	// Compilar temporariamente a CLI
	t.Log("Compilando CLI para E2E...")
	cmd := exec.Command("go", "build", "-o", "crom-agente-cli-e2e", "./cmd/crom-agente-cli")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Falha ao compilar CLI: %v", err)
	}
	defer exec.Command("rm", "-f", "crom-agente-cli-e2e").Run()

	t.Log("Executando crom-agente-cli-e2e diag...")
	runCmd := exec.Command("./crom-agente-cli-e2e", "diag")
	var out bytes.Buffer
	runCmd.Stdout = &out
	runCmd.Stderr = &out

	if err := runCmd.Run(); err != nil {
		t.Fatalf("Falha ao rodar diag: %v, output: %s", err, out.String())
	}

	output := out.String()
	if !strings.Contains(output, "Diagnóstico do Ambiente crom-agente") {
		t.Errorf("Saída E2E não contém o banner esperado: %s", output)
	}
	if !strings.Contains(output, "git") {
		t.Errorf("Saída E2E não contém checagem do git: %s", output)
	}
}
