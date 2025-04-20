# Introduction

The microservice core is a small library which contains helper functions and basic setups for golang microservices. It's just a minimum set for messaging, rest and some helpers, but not so mightful as [Goa](https://github.com/goadesign/goa)

In the near future goa could be a better option to replace this. 

# Usage

## How to use config?
1. define config struct with all required values (include BaseConfig too!)
````go
type exampleConfig struct {
    core.BaseConfig `mapstructure:",squash"`
    TestValue       string `mapstructure:"testValue"`
}
````

you can also define extra structs inside the config struct:
````go
type exampleConfig struct {
    core.BaseConfig `mapstructure:",squash"`
    TestValue       string `mapstructure:"testValue"`
    OAuth           struct {
        ServerUrl    string `mapstructure:"serverUrl"`
        ClientId     string `mapstructure:"clientId"`
        ClientSecret string `mapstructure:"clientSecret"`
    } `mapstructure:"oAuth"`
}
````

2. Create global config struct instance and load config with LoadConfig() function. Provide Prefix for ENV variables and map with default values
````go
var ExampleConfig exampleConfig

func LoadConfig() error {
	err := core.LoadConfig("EXAMPLE", &ExampleConfig, getDefaults())
	if err != nil {
		return err
	}
	return nil
}

func getDefaults() map[string]any {
    return map[string]any{
    "testValue": "ciao",
    }
}
````

Note: The baseconfig can be also used by using envconfig. In this case the envconfig package is required and a envconfig processing before the start. But be aware that both variants in the same time can have interferences in variable loading (envconfig default can override viper loading and the other way arround). 