package handlerruntime

type Data struct {
	ID    string         `mapstructure:"id"`
	Name  string         `mapstructure:"name"`
	Other map[string]any `mapstructure:",remain"`
}

type Resource struct {
	Type string `mapstructure:"type"`
	Data Data   `mapstructure:"data"`
}

type LoadResourceResponse struct {
	Resources []Resource `mapstructure:"resources"`

	Tasks []struct {
		Name string         `mapstructure:"name"`
		Ctx  map[string]any `mapstructure:"ctx"`
	} `mapstructure:"tasks"`
}
