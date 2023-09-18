package config

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/melisource/fury_go-dev-base-3-v2/cmd/api/domain"

	furyConfig "github.com/mercadolibre/fury_go-toolkit-config/v2/pkg/config"
	"go.uber.org/zap/zapcore"

	"github.com/mercadolibre/fury_go-core/pkg/log"

	"github.com/spf13/viper"
)

const (
	staticScope       = "SCOPE"
	staticFuryProfile = "default"
	envLocal          = "local"
	envSandbox        = "sandbox"
	envProduction     = "production"
	fileConfDefault   = "/app/resources"
	fileExtension     = "yaml"
	LimitSearchRows   = 100

	fileNotFoundError = "no such file or directory"
)

var (
	mapLocalFile = map[string]string{
		envLocal:      fileConfDefault,
		envSandbox:    fileConfDefault,
		envProduction: fileConfDefault,
	}
)

type Environment struct {
	LogLevel       zapcore.Level `mapstructure:"logLevel"`
	ScopeContainer string
	MySQLConfig    domain.MySQL `mapstructure:"mysqlconfig"`
}

type ConnectionConfig struct {
	TimeoutSeconds         time.Duration
	DisableTimeout         bool
	MaxIdleConnections     int
	MaxRetries             int
	RetryDelayMilliSeconds time.Duration
}

func InitConfig() Environment {
	instance := Environment{}

	time.Local = time.UTC

	handlerEnvViper(&instance)
	instance.DecryptEnvironment()
	return instance
}

func (e *Environment) DecryptEnvironment() {
	if viper.GetString(staticScope) == envLocal {
		e.MySQLConfig.Password = os.Getenv("MYSQL_PASSWORD")
	}
}

func handlerEnvViper(environment *Environment) {
	viper.AutomaticEnv()

	viper.SetDefault(staticScope, envLocal)
	environment.ScopeContainer = handlerFormatScope()

	workingdir, err := os.Getwd()
	if err != nil {
		log.Panic(context.Background(), err.Error())
	}
	dir := strings.Split(workingdir, "/cmd/api")[0]

	viper.SetConfigName(environment.ScopeContainer)
	viper.SetConfigType(fileExtension)
	viper.AddConfigPath(mapLocalFile[environment.ScopeContainer])
	viper.AddConfigPath(dir + "/resources")
	viper.AutomaticEnv()

	loadEnvironmentFile()
	loadFuryConfig()
	loadEnvironment(environment)
}

func loadEnvironment(environment *Environment) {
	for _, k := range viper.AllKeys() {
		value := viper.Get(k)
		if _, ok := value.(string); ok {
			viper.Set(k, os.ExpandEnv(viper.GetString(k)))
		}
	}

	if err := viper.Unmarshal(&environment); err != nil {
		msg := fmt.Sprintf("[event: fault_conf_init][service: unmarshal yaml] Could not unmarshal configs %s", err)
		log.Panic(context.Background(), msg)
	}
}

func loadEnvironmentFile() {
	if err := viper.ReadInConfig(); err != nil {
		msg := fmt.Sprintf("[event: fault_conf_init][service: read yaml] Could not loading configs %s", err)
		log.Panic(context.Background(), msg)
	}

}

func loadFuryConfig() {
	configData, errLoadConfig := furyConfig.Read(staticFuryProfile)
	if errLoadConfig != nil {
		if strings.Contains(errLoadConfig.Error(), fileNotFoundError) {
			return
		}
		msg := fmt.Sprintf("[event: fault_conf_init][service: unmarshal yaml] Could not load fury configs %s", errLoadConfig)
		log.Panic(context.Background(), msg)
	}

	if err := viper.MergeConfig(bytes.NewBuffer(configData)); err != nil {
		msg := fmt.Sprintf("[event: fault_conf_init][service: read fury config] Could not loading fury configs %s", err)
		log.Panic(context.Background(), msg)
	}
}

func handlerFormatScope() string {
	scope := viper.GetString(staticScope)

	if strings.Contains(strings.ToLower(scope), envSandbox) || strings.Contains(strings.ToLower(scope), "sbx") {
		return envSandbox
	}

	if strings.Contains(strings.ToLower(scope), envProduction) || strings.Contains(strings.ToLower(scope), "prd") {
		return envProduction
	}

	return envLocal
}

func (e *Environment) IsProduction() bool {
	return e.ScopeContainer == envProduction
}
