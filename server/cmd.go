package server

import (
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/application"
	"github.com/spf13/cobra"
)

func Execute() {
	cobra.OnInitialize(func() {
		// 初始化ioc
		cobra.CheckErr(ioc.ConfigIocObject(configReq))

		// 根据配置文件，补充信息到 cobra root command
		Root.Use = application.Get().ApplicationName()
		Root.Short = application.Get().AppDescription
		Root.Long = application.Get().AppDescription
	})

	// 从root 进行启动
	cobra.CheckErr(Root.Execute())
}

var (
	confType string
	confFile string
)

var Root = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动服务",
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(RunServ(cmd.Context()))
	},
}

func init() {
	Root.PersistentFlags().StringVarP(&confType, "config-type", "t", "file", "the service config type [file/env]")
	Root.PersistentFlags().StringVarP(&confFile, "config-file", "f", "etc/application.yaml", "the service config from file")
	Root.AddCommand(startCmd)
}
