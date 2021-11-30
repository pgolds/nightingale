package models

type NotifyTemplate struct {
	IsAlert bool
	IsMachineDep bool
	Sname string
	Ident string
	Classpath string
	Metric string
	Tags string
	Value string
	Status string
	ReadableExpression string
	TriggerTime string
	RuleId int64
	EventId int64
}
