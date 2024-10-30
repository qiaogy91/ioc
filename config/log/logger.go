package log

import (
	"context"
	"github.com/qiaogy91/ioc"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	ins = &Logger{
		lock:       new(sync.Mutex),
		subLoggers: make(map[string]*slog.Logger),
		Trace:      false,
		Level:      slog.LevelDebug,
		Filename:   "logs/app.log",
		MaxSize:    10,
		MaxAge:     20,
		MaxBackups: 6,
		LocalTime:  true,
		Compress:   true,
		Deep:       3,
	}
)

type Logger struct {
	ioc.ObjectImpl
	lock       sync.Locker
	subLoggers map[string]*slog.Logger
	root       *slog.Logger
	Trace      bool       `json:"trace" yaml:"trace"`           // 开启trace
	Level      slog.Level `json:"level" yaml:"level"`           // 级别
	Filename   string     `json:"filename" yaml:"filename"`     // 日志文件名
	MaxSize    int        `json:"maxSize" yaml:"maxSize"`       // 触发滚动的文件大小，单位 megabytes
	MaxAge     int        `json:"maxAge" yaml:"maxAge"`         // 触发滚动的文件保留时间，单位 day
	MaxBackups int        `json:"maxBackups" yaml:"maxBackups"` // 滚动时，旧文件保留多少分
	LocalTime  bool       `json:"LocalTime" yaml:"LocalTime"`   // 备份文件的时间格式
	Compress   bool       `json:"compress" yaml:"compress"`     // 备份文件是否进行压缩
	Deep       int        `json:"deep" yaml:"deep"`             // 文件路径深度

}

func (l *Logger) Name() string  { return AppName }
func (l *Logger) Priority() int { return 101 }
func (l *Logger) Close(ctx context.Context) error {
	selfLog := l.SubLogger(AppName)
	selfLog.Debug("closed completed", slog.String("namespace", ioc.ConfigNamespace))
	return nil
}

func (l *Logger) replaceFilePath(filePath string, deep int) string {
	parts := strings.Split(filepath.ToSlash(filePath), "/") // 将路径分割为目录部分
	// 少于指定层级的，直接返回原路径
	if len(parts) <= deep {
		return filePath
	}
	// 保留最后 n 级目录
	parts = parts[len(parts)-deep:]
	return filepath.Join(parts...)
}
func (l *Logger) handlerOpts() *slog.HandlerOptions {
	return &slog.HandlerOptions{
		AddSource: true,
		Level:     l.Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				if src, ok := a.Value.Any().(*slog.Source); ok {
					src.File = l.replaceFilePath(src.File, l.Deep)
					return slog.Attr{Key: slog.SourceKey, Value: slog.AnyValue(src)}
				}
			}
			return a
		},
	}
}
func (l *Logger) HandlerFile() *slog.JSONHandler {
	file := &lumberjack.Logger{
		Filename:   l.Filename,
		MaxSize:    l.MaxSize,
		MaxAge:     l.MaxAge,
		MaxBackups: l.MaxBackups,
		LocalTime:  true,
		Compress:   true,
	}

	return slog.NewJSONHandler(file, l.handlerOpts())
}
func (l *Logger) HandlerConsole() *slog.TextHandler {
	return slog.NewTextHandler(os.Stdout, l.handlerOpts())
}

func (l *Logger) SubLogger(name string) *slog.Logger {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.root == nil {
		return nil
	}

	if _, ok := l.subLoggers[name]; !ok {
		l.subLoggers[name] = l.root.With(slog.String(SubLoggerKey, name))
	}

	return l.subLoggers[name]
}

func (l *Logger) Init() {
	handlers := &MultiHandler{
		hs: []slog.Handler{
			l.HandlerConsole(),
			l.HandlerFile(),
		},
	}

	// 如果开启追踪，则增加一个handler 用来将日志发送到 Otlp
	if l.Trace {
		ioc.OtlpMustEnabled()
		handlers.hs = append(handlers.hs, otelslog.NewHandler("trace-handler"))
	}

	l.root = slog.New(handlers)
}

func init() {
	ioc.Config().Registry(ins)
}
