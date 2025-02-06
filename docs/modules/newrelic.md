# New Relic 
The New Relic module allows the Chaki application to export telemetry data to New Relic, enabling detailed monitoring and inspection. This module includes various submodules, such as `server`, `client`, `kafka`, and `postgresql`, which can be used to capture specific telemetry data.

To use the New Relic module in your application, you need to declare it as you would with any other module. However, the New Relic module itself requires the use of these submodules to function effectively.

### Usage

First, import the necessary New Relic submodules:

```go
import (
    nrorm "github.com/Trendyol/chaki/modules/newrelic/orm"
    nrserver "github.com/Trendyol/chaki/modules/newrelic/server"
)

// ...

app := chaki.New()

app.Use(
// Other modules...

newrelic.Module(
    nrserver.WithServer(),
    nrorm.WithPostgresql(),
)

// ...
```

### Config

The New Relic configuration requires the following parameters:
```yaml
newrelic:
  enabled: true
  logenabled: true
  appname: "server-with-newrelic"
  license: "top-secret-license-key"
```
- enabled: Enable or disable the New Relic agent.
- logenabled: enables or disables all application logging features and data collection
- appname: The name of your application as it will appear in New Relic.
- licensekey: Your New Relic license key for authentication.


You can pass more specific New Relic options to the module through its constructor using configuration options. 


```go
app.Use(
  // other modules
  newrelic.Module(
		WithNewrelicConfigOptions(newrelic.ConfigAIMonitoringEnabled(true)),
  )
)

```