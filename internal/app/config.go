package app

import (
	"github.com/core-go/core"
	"github.com/core-go/core/header"
	"github.com/core-go/core/server"
	mid "github.com/core-go/log/middleware"
	"github.com/core-go/log/zap"
)

type Config struct {
	Server     server.ServerConfig `mapstructure:"server"`
	Mongo      MongoConfig         `mapstructure:"mongo"`
	Log        log.Config          `mapstructure:"log"`
	Response   header.Config       `mapstructure:"response"`
	MiddleWare mid.LogConfig       `mapstructure:"middleware"`
	Action     *core.ActionConfig  `mapstructure:"action"`
}

type MongoConfig struct {
	Uri      string `yaml:"uri" mapstructure:"uri" json:"uri,omitempty" gorm:"column:uri" bson:"uri,omitempty" dynamodbav:"uri,omitempty" firestore:"uri,omitempty"`
	Database string `yaml:"database" mapstructure:"database" json:"database,omitempty" gorm:"column:database" bson:"database,omitempty" dynamodbav:"database,omitempty" firestore:"database,omitempty"`
}
