package handlerruntime

type Data struct {
	ID    string                 `mapstructure:"id"`
	Name  string                 `mapstructure:"name"`
	Other map[string]interface{} `mapstructure:",remain"`
}

type Resource struct {
	Type string `mapstructure:"type"`
	Data Data   `mapstructure:"data"`
}

type LoadResourceResponse struct {
	Resources []Resource `mapstructure:"resources"`

	PendingTasks []struct {
		Name string      `mapstructure:"name"`
		Ctx  interface{} `mapstructure:"ctx"`
	} `mapstructure:"pendingTasks"`
}
