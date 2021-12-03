你收到一条告警信息,策略名称为{{.Sname}},
{{if .IsMachineDep}}告警设备为{.Ident}},
挂载节点为{{.Classpath}}{{end}},
当前值为{{.Value}},
触发时间为{{.TriggerTime}},
请尽快处理。