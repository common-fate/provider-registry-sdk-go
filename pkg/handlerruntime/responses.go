package handlerruntime

type Data struct {
	ID    string         `mapstructure:"id" json:"id"`
	Name  string         `mapstructure:"name" json:"name"`
	Other map[string]any `mapstructure:",remain"`
}

type Resource struct {
	Type string `mapstructure:"type" json:"type"`
	Data Data   `mapstructure:"data" json:"data"`
}

type LoadResourceResponse struct {
	Resources []Resource         `mapstructure:"resources"`
	Tasks     []LoadResourceTask `mapstructure:"tasks" json:"tasks"`
}

type LoadResourceTask struct {
	Name string         `mapstructure:"name" json:"name"`
	Ctx  map[string]any `mapstructure:"ctx" json:"ctx"`
}
