级别状态：{{.Status}}
策略名称：{{.Sname}}
{{if .IsMachineDep}}告警设备：{{.Ident}}
挂载节点：{{.Classpath}}{{end}}
指标标签：{{.Tags}}
当前值：{{.Value}}
报警说明：{{.ReadableExpression}}
触发时间：{{.TriggerTime}}
