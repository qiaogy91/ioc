package server

import (
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/application"
	"github.com/spf13/cobra"
)

var (
	confType string
	confFile string
)

var Root = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Start() {
	Root.AddCommand(startCmd)
	Execute()
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动服务",
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(RunServ(cmd.Context()))
	},
}

func Execute() {
	cobra.OnInitialize(func() {
		switch confType {
		case "file":
			configReq.ConfigFile.Enabled = true
			configReq.ConfigFile.Path = confFile
		default:
			configReq.ConfigEnv.Enabled = true
		}

		// 初始化ioc
		cobra.CheckErr(ioc.ConfigIocObject(configReq))

		// 补充Root命令信息
		Root.Use = application.Get().AppName
		Root.Short = application.Get().AppDescription
		Root.Long = application.Get().AppDescription
	})
	cobra.CheckErr(Root.Execute())
}

func init() {
	Root.PersistentFlags().StringVarP(&confType, "config-type", "t", "file", "the service config type [file/env]")
	Root.PersistentFlags().StringVarP(&confFile, "config-file", "f", "etc/application.yaml", "the service config from file")
}
