// Copyright 2020 Layer5, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/api/grpc"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
	internalconfig "github.com/layer5io/meshery-osm/internal/config"
	"github.com/layer5io/meshery-osm/osm"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils"
	"os"
	"path"
	"time"
)

const (
	serviceName = "osm-adapter"
)

func main() {
	log, err := logger.New(serviceName, logger.Options{Format: logger.JsonLogFormat, DebugLevel: false})
	if err != nil {
		fmt.Println("Logger Init Failed", err.Error())
		os.Exit(1)
	}

	cfg, err := internalconfig.New(configprovider.ViperKey)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	kubeconfigHandler, err := internalconfig.NewKubeconfigBuilder(configprovider.ViperKey)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	service := &grpc.Service{}
	_ = cfg.GetObject(adapter.ServerKey, &service)

	service.Handler = osm.New(cfg, log, kubeconfigHandler)
	service.Channel = make(chan interface{}, 100)
	service.StartedAt = time.Now()
	err = grpc.Start(service, nil)
	if err != nil {
		log.Error(grpc.ErrGrpcServer(err))
		os.Exit(1)
	}
}

// This init function can help adapters create the configuration logic work well, so do not remove it although that's
// not a good idea.
func init() {
	err := os.MkdirAll(path.Join(utils.GetHome(), ".meshery", "bin"), 0750)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
