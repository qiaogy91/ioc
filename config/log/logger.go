package log

import (
	"fmt"
	"github.com/qiaogy91/ioc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	ioc.ObjectImpl
	root       *zerolog.Logger
	lock       sync.Mutex
	loggers    map[string]*zerolog.Logger
	CallerDeep int           `json:"callerDeep" yaml:"callerDeep"` // 0 为打印日志全路径, 默认打印2层路径
	Level      zerolog.Level `json:"level" yaml:"level"`           // 日志的级别, 默认Debug
	NoColor    bool          `json:"noColor" yaml:"noColor"`       // 控制台输出颜色
	FilePath   string        `json:"filePath" yaml:"filePath"`     // 文件路径，配置后则认为开启了文件日志
	MaxSize    int           `json:"maxSize" yaml:"maxSize"`       // 单位M，默认100M
	MaxBackups int           `json:"maxBackups" yaml:"maxBackups"` // 默认保存6个
	MaxAge     int           `json:"maxAge" yaml:"maxAge"`         // 保存多久
	Compress   bool          `json:"compress" yaml:"compress"`     // 是否压缩
}

func (l *Logger) ConsoleWriter() io.Writer {
	output := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.NoColor = l.NoColor
		w.TimeFormat = time.RFC3339
	})

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%-6s", i))
	} // 日志级别为宽度6个字符的字串
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	} // 消息内容以字符串方式输出
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	} // 字段名称后加双引号
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	} // 字段值转为大写
	return output
}

func (l *Logger) FileWriter() io.Writer {
	return &lumberjack.Logger{
		Filename:   l.FilePath,
		MaxSize:    l.MaxSize,
		MaxAge:     l.MaxAge,
		MaxBackups: l.MaxBackups,
		Compress:   l.Compress,
	}
}

func (l *Logger) Name() string  { return AppName }
func (l *Logger) Priority() int { return 102 }
func (l *Logger) Init() {
	var writers []io.Writer
	writers = append(writers, l.ConsoleWriter())
	if l.FilePath != "" {
		writers = append(writers, l.FileWriter())
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	root := zerolog.New(io.MultiWriter(writers...)).With().Timestamp()
	if l.CallerDeep > 0 {
		root = root.Caller()
		zerolog.CallerMarshalFunc = l.CallerMarshalFunc
	}

	l.SetRoot(root.Logger().Level(l.Level))
}
func (l *Logger) CallerMarshalFunc(pc uintptr, file string, line int) string {
	if l.CallerDeep == 0 {
		return file
	}

	short := file
	count := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			count++
		}

		if count >= l.CallerDeep {
			break
		}
	}
	file = short
	info := file + ":" + strconv.Itoa(line)
	if len(info) < 25 {
		info = fmt.Sprintf("%-*s", 25, info)
	}
	return info
}
func (l *Logger) SetRoot(r zerolog.Logger) {
	l.root = &r
}
func (l *Logger) Logger(name string) *zerolog.Logger {
	l.lock.Lock()
	defer l.lock.Unlock()

	if _, ok := l.loggers[name]; !ok {
		sub := l.root.With().Str(SubLoggerKey, name).Logger()
		l.loggers[name] = &sub
	}

	return l.loggers[name]
}

func init() {
	ioc.Config().Registry(&Logger{
		loggers: make(map[string]*zerolog.Logger),
	})
}
