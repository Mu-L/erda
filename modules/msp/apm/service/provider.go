// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	logs "github.com/erda-project/erda-infra/base/logs"
	servicehub "github.com/erda-project/erda-infra/base/servicehub"
	transport "github.com/erda-project/erda-infra/pkg/transport"
	metricpb "github.com/erda-project/erda-proto-go/core/monitor/metric/pb"
	pb "github.com/erda-project/erda-proto-go/msp/apm/service/pb"
	"github.com/erda-project/erda/pkg/common/apis"
	perm "github.com/erda-project/erda/pkg/common/permission"
)

type View struct {
	ViewType string   `file:"type"`
	Charts   []string `file:"charts"`
}

type config struct {
	View []*View
}

// +provider
type provider struct {
	Cfg               *config
	Log               logs.Logger
	Register          transport.Register
	apmServiceService *apmServiceService
	Metric            metricpb.MetricServiceServer `autowired:"erda.core.monitor.metric.MetricService"`
	Perm              perm.Interface               `autowired:"permission"`
}

func GetView(c *config, key string) *View {
	for _, v := range c.View {
		if v.ViewType == key {
			return v
		}
	}
	return nil
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.apmServiceService = &apmServiceService{p}
	if p.Register != nil {
		pb.RegisterApmServiceServiceImp(p.Register, p.apmServiceService, apis.Options())
	}
	return nil
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	switch {
	case ctx.Service() == "erda.msp.apm.service.ApmServiceService" || ctx.Type() == pb.ApmServiceServiceServerType() || ctx.Type() == pb.ApmServiceServiceHandlerType():
		return p.apmServiceService
	}
	return p
}

func init() {
	servicehub.Register("erda.msp.apm.service", &servicehub.Spec{
		Services:             pb.ServiceNames(),
		Types:                pb.Types(),
		OptionalDependencies: []string{"service-register"},
		Description:          "",
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}