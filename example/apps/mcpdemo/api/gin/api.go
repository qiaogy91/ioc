package gin

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	mcpcore "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/qiaogy91/ioc"
	iocgin "github.com/qiaogy91/ioc/config/gin"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/qiaogy91/ioc/example/apps/mcpdemo"
)

// Handler MCP 工具注册器
type Handler struct {
	ioc.ObjectImpl
	log       *slog.Logger
	svc       mcpdemo.Controller
	mcpServer *server.MCPServer

	// 配置项
	Stateless bool   `json:"stateless" yaml:"stateless"`
	Transport string `json:"transport" yaml:"transport"`
}

func (h *Handler) Name() string  { return mcpdemo.HandlerName }
func (h *Handler) Priority() int { return 411 }

func (h *Handler) Init() {
	h.log = log.Sub(mcpdemo.HandlerName)
	h.svc = mcpdemo.GetController()

	// 创建 MCP Server
	h.mcpServer = server.NewMCPServer("mcpdemo", "1.0.0")

	// 注册 MCP 工具
	h.registerTools()

	// 创建 StreamableHTTP transport
	var opts []server.StreamableHTTPOption
	if h.Stateless {
		opts = append(opts, server.WithStateLess(true))
	}
	streamableServer := server.NewStreamableHTTPServer(h.mcpServer, opts...)

	// 挂载到 gin 路由
	// ModuleRouter 会根据 Name() 创建 /mcpdemo_mcp 前缀
	// 最终路径为 /mcpdemo_mcp/mcp
	router := iocgin.ModuleRouter(h)
	router.Any("/mcp", gin.WrapH(streamableServer))
}

// registerTools 注册 MCP 工具
func (h *Handler) registerTools() {
	// 注册 get_user 工具
	getUserTool := mcpcore.Tool{
		Name:        "get_user",
		Description: "Get user information by ID",
		InputSchema: mcpcore.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "integer",
					"description": "User ID",
				},
			},
			Required: []string{"id"},
		},
	}
	h.mcpServer.AddTool(getUserTool, h.handleGetUser)
}

// handleGetUser 处理 get_user 工具调用
func (h *Handler) handleGetUser(ctx context.Context, request mcpcore.CallToolRequest) (*mcpcore.CallToolResult, error) {
	args := request.GetArguments()
	id := int(args["id"].(float64))

	user, err := h.svc.GetUser(ctx, id)
	if err != nil {
		return &mcpcore.CallToolResult{
			IsError: true,
			Content: []mcpcore.Content{
				mcpcore.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error: %v", err),
				},
			},
		}, nil
	}

	return &mcpcore.CallToolResult{
		Content: []mcpcore.Content{
			mcpcore.TextContent{
				Type: "text",
				Text: fmt.Sprintf("User: ID=%d, Name=%s, Age=%d", user.ID, user.Name, user.Age),
			},
		},
	}, nil
}

func (h *Handler) Close(ctx context.Context) error {
	if h.log != nil {
		h.log.Debug("closed completed", slog.String("namespace", ioc.ApiNamespace))
	}
	return nil
}

func init() {
	ioc.Api().Registry(&Handler{})
}
