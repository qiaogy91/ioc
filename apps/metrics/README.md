
#### 说明
对于api metric 指标分别有两种实现类型（opentelemetry、prometheus），每种类型有两种适用框架（gin、go-restful）
同时只能使用其中某一种框架的某一个实现（对restful 选择opentelemetry 实现）
```shell
➜  ioc git:(main) ✗ tree apps/metrics -L 2
apps/metrics
├── README.md.md
├── otlp
│   ├── gin
│   └── restful
└── prom
    ├── gin
    └── restful
```