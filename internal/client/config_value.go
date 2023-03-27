package client

type ConfigValue struct {
	Id                 int64 `json:"id"`
	EnvironmentId      int64 `json:"environment_id"`
	ConfigurationKeyId int64 `json:"configuration_key_id"`

	Key        string  `json:"key"`
	ValueType  string  `json:"value_type"`
	StrValue   string  `json:"str_value"`
	IntValue   int64   `json:"int_value"`
	FloatValue float64 `json:"float_value"`
	BoolValue  bool    `json:"bool_value"`
}
