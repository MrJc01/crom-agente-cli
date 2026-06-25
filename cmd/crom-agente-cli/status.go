package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// Structs de telemetria locais para evitar problemas com pacotes internal do Go

type MCPServerStatus struct {
	Name      string   `json:"name"`
	Mode      string   `json:"mode"`
	ToolCount int      `json:"tool_count"`
	Tools     []string `json:"tools"`
	Running   bool     `json:"running"`
}

type TerminalTelemetry struct {
	ID        string    `json:"id"`
	PID       int       `json:"pid"`
	Name      string    `json:"name"`
	Closed    bool      `json:"closed"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProcessTelemetry struct {
	ID           string    `json:"id"`
	Command      string    `json:"command"`
	PID          int       `json:"pid"`
	Status       string    `json:"status"`
	StartedAt    time.Time `json:"started_at"`
	IsBackground bool      `json:"is_background"`
}

type TaskItem struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
}

type AgentState struct {
	ID                      string              `json:"id,omitempty"`
	Name                    string              `json:"name,omitempty"`
	Status                  string              `json:"status,omitempty"`
	DiretorioAtual          string              `json:"diretorio_atual"`
	ArquivosFocados         []string            `json:"arquivos_focados"`
	TarefaEmAndamento       string              `json:"tarefa_em_andamento"`
	UltimoStatus            string              `json:"ultimo_status"`
	StatusOperacional       string              `json:"status_operacional"`
	ModoCognitivo           string              `json:"modo_cognitivo"`
	LogsRelevantes          []string            `json:"logs_relevantes"`
	TokensGastos            int                 `json:"tokens_gastos"`
	TotalTurnos             int                 `json:"total_turnos"`
	Timestamp               time.Time           `json:"timestamp"`
	Plan                    []TaskItem          `json:"plan,omitempty"`
	BrowserURL              string              `json:"browser_url,omitempty"`
	FilesCreated            int                 `json:"files_created"`
	FilesValidated          int                 `json:"files_validated"`
	ToolCallsEmitted        int                 `json:"tool_calls_emitted"`
	ToolCallsFromTextParse  int                 `json:"tool_calls_from_text_parse"`
	CircuitBreakerTriggered bool                `json:"circuit_breaker_triggered"`
	ActiveTerminals         []TerminalTelemetry `json:"active_terminals"`
	ActiveProcesses         []ProcessTelemetry  `json:"active_processes"`
	CurrentStep             string              `json:"current_step"`
	CurrentStepDurationMs   int64               `json:"current_step_duration_ms"`
}

type BrowserTelemetry struct {
	Active bool   `json:"active"`
	URL    string `json:"url"`
}

type AgentTelemetry struct {
	WorkspaceName string            `json:"workspace_name"`
	IsRunning     bool              `json:"is_running"`
	AgentState    AgentState        `json:"agent_state"`
	Browser       BrowserTelemetry  `json:"browser"`
	MCPServers    []MCPServerStatus `json:"mcp_servers"`
}

var (
	// Definição de estilos Lipgloss
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#A3E635")).
			Background(lipgloss.Color("#1E293B")).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#38BDF8")).
			MarginTop(1)

	sectionStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#475569")).
			Padding(0, 1).
			Width(78)

	greenText = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ADE80"))
	redText   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F87171"))
	blueText  = lipgloss.NewStyle().Foreground(lipgloss.Color("#60A5FA"))
	grayText  = lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8"))
	yellowText = lipgloss.NewStyle().Foreground(lipgloss.Color("#FBBF24"))
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Mostra a telemetria em tempo real e estado do agente no workspace (Stage 9)",
	Long:  `O comando status consulta o Daemon local e exibe informações sobre processos, terminais, navegador e estado cognitivo do agente.`,
	Run: func(cmd *cobra.Command, args []string) {
		absPath, err := filepath.Abs(workspacePath)
		if err != nil {
			absPath = workspacePath
		}
		
		fmt.Printf("🔍 Consultando telemetria para o workspace: %s...\n", absPath)
		
		client := &http.Client{Timeout: 3 * time.Second}
		reqURL := fmt.Sprintf("http://127.0.0.1:9090/api/agent/telemetry?workspace=%s", url.QueryEscape(absPath))
		
		resp, err := client.Get(reqURL)
		if err != nil {
			fmt.Println(redText.Render(fmt.Sprintf("❌ Erro ao conectar ao Daemon na porta 9090: %v", err)))
			fmt.Println(grayText.Render("Certifique-se de que o daemon persistente do crom-agente está em execução."))
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println(redText.Render(fmt.Sprintf("❌ Daemon retornou erro HTTP %d ao consultar telemetria.", resp.StatusCode)))
			os.Exit(1)
		}

		var telemetry AgentTelemetry
		if err := json.NewDecoder(resp.Body).Decode(&telemetry); err != nil {
			fmt.Println(redText.Render(fmt.Sprintf("❌ Falha ao parsear telemetria do Daemon: %v", err)))
			os.Exit(1)
		}

		renderTelemetry(telemetry)
	},
}

func renderTelemetry(t AgentTelemetry) {
	fmt.Println()
	fmt.Println(titleStyle.Render(fmt.Sprintf("🛰️  TELEMETRIA: %s", t.WorkspaceName)))
	
	// Bloco 1: Status Geral
	runningStr := redText.Render("Inativo (Idle)")
	if t.IsRunning {
		runningStr = greenText.Render("Ativo (Running)")
	}
	
	stepStr := t.AgentState.CurrentStep
	if stepStr == "" {
		stepStr = "-"
	}
	
	durationStr := "-"
	if t.AgentState.CurrentStepDurationMs > 0 {
		durationStr = fmt.Sprintf("%.2fs", float64(t.AgentState.CurrentStepDurationMs)/1000.0)
	}

	generalInfo := fmt.Sprintf(
		"Loop: %s\nOperação: %s\nModo Cognitivo: %s\nEtapa Atual: %s (%s)\nTokens Gastos: %s | Turnos: %d",
		runningStr,
		blueText.Render(t.AgentState.StatusOperacional),
		yellowText.Render(t.AgentState.ModoCognitivo),
		blueText.Render(stepStr),
		grayText.Render(durationStr),
		greenText.Render(fmt.Sprintf("%d", t.AgentState.TokensGastos)),
		t.AgentState.TotalTurnos,
	)
	
	fmt.Println(headerStyle.Render("📊 Status Geral"))
	fmt.Println(sectionStyle.Render(generalInfo))

	// Bloco 2: Terminais e Processos
	termInfo := ""
	activeTermsCount := 0
	for _, term := range t.AgentState.ActiveTerminals {
		if !term.Closed {
			activeTermsCount++
			termInfo += fmt.Sprintf("💻 Terminal %s (PID: %d) - %s\n", blueText.Render(term.ID), term.PID, greenText.Render("Aberto"))
		}
	}
	if activeTermsCount == 0 {
		termInfo += grayText.Render("Nenhum terminal ativo.\n")
	}
	
	procInfo := ""
	activeProcsCount := 0
	for _, proc := range t.AgentState.ActiveProcesses {
		if proc.Status == "running" {
			activeProcsCount++
			bgStr := "Foreground"
			if proc.IsBackground {
				bgStr = "Background"
			}
			statusColor := yellowText
			procInfo += fmt.Sprintf("⚙️  %s (PID: %d) | %s | %s\n   %s\n",
				blueText.Render(proc.ID),
				proc.PID,
				bgStr,
				statusColor.Render(proc.Status),
				grayText.Render(proc.Command),
			)
		}
	}
	if activeProcsCount == 0 {
		procInfo += grayText.Render("Nenhum subprocesso ativo.\n")
	}

	fmt.Println(headerStyle.Render("💻 Terminais e Subprocessos Ativos"))
	fmt.Println(sectionStyle.Render(strings.TrimSpace(termInfo + "\n" + procInfo)))

	// Bloco 3: Browser & MCP
	browserStr := grayText.Render("Navegador Fechado")
	if t.Browser.Active {
		browserStr = fmt.Sprintf("%s\nURL: %s", greenText.Render("Navegador Ativo"), blueText.Render(t.Browser.URL))
	}
	
	mcpStr := ""
	if len(t.MCPServers) == 0 {
		mcpStr = grayText.Render("Nenhum servidor MCP global ativo.")
	} else {
		for _, mcp := range t.MCPServers {
			status := redText.Render("offline")
			if mcp.Running {
				status = greenText.Render("online")
			}
			mcpStr += fmt.Sprintf("🔌 %s (%s) | %s | %d ferramentas\n",
				blueText.Render(mcp.Name),
				mcp.Mode,
				status,
				mcp.ToolCount,
			)
		}
	}

	fmt.Println(headerStyle.Render("🌐 Navegador & Integrações MCP"))
	fmt.Println(sectionStyle.Render(browserStr + "\n\n" + strings.TrimSpace(mcpStr)))

	// Bloco 4: Logs Recentes
	logsStr := ""
	if len(t.AgentState.LogsRelevantes) == 0 {
		logsStr = grayText.Render("Nenhum log registrado.")
	} else {
		for i, logLine := range t.AgentState.LogsRelevantes {
			logsStr += fmt.Sprintf("[%d] %s\n", i+1, grayText.Render(logLine))
		}
	}
	fmt.Println(headerStyle.Render("📝 Histórico de Logs Recentes"))
	fmt.Println(sectionStyle.Render(strings.TrimSpace(logsStr)))
	fmt.Println()
}
