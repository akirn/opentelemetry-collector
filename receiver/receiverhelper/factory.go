// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package receiverhelper

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
)

// FactoryOption apply changes to ReceiverOptions.
type FactoryOption func(o *factory)

// WithTraces overrides the default "error not supported" implementation for CreateTraceReceiver.
func WithTraces(createTraceReceiver CreateTraceReceiver) FactoryOption {
	return func(o *factory) {
		o.createTraceReceiver = createTraceReceiver
	}
}

// WithMetrics overrides the default "error not supported" implementation for CreateMetricsReceiver.
func WithMetrics(createMetricsReceiver CreateMetricsReceiver) FactoryOption {
	return func(o *factory) {
		o.createMetricsReceiver = createMetricsReceiver
	}
}

// WithLogs overrides the default "error not supported" implementation for CreateLogsReceiver.
func WithLogs(createLogsReceiver CreateLogsReceiver) FactoryOption {
	return func(o *factory) {
		o.createLogsReceiver = createLogsReceiver
	}
}

// CreateDefaultConfig is the equivalent of component.ReceiverFactory.CreateDefaultConfig()
type CreateDefaultConfig func() config.Receiver

// CreateTraceReceiver is the equivalent of component.ReceiverFactory.CreateTracesReceiver()
type CreateTraceReceiver func(context.Context, component.ReceiverCreateParams, config.Receiver, consumer.Traces) (component.TracesReceiver, error)

// CreateMetricsReceiver is the equivalent of component.ReceiverFactory.CreateMetricsReceiver()
type CreateMetricsReceiver func(context.Context, component.ReceiverCreateParams, config.Receiver, consumer.Metrics) (component.MetricsReceiver, error)

// CreateLogsReceiver is the equivalent of component.ReceiverFactory.CreateLogsReceiver()
type CreateLogsReceiver func(context.Context, component.ReceiverCreateParams, config.Receiver, consumer.Logs) (component.LogsReceiver, error)

type factory struct {
	cfgType               config.Type
	createDefaultConfig   CreateDefaultConfig
	createTraceReceiver   CreateTraceReceiver
	createMetricsReceiver CreateMetricsReceiver
	createLogsReceiver    CreateLogsReceiver
}

// NewFactory returns a component.ReceiverFactory.
func NewFactory(
	cfgType config.Type,
	createDefaultConfig CreateDefaultConfig,
	options ...FactoryOption) component.ReceiverFactory {
	f := &factory{
		cfgType:             cfgType,
		createDefaultConfig: createDefaultConfig,
	}
	for _, opt := range options {
		opt(f)
	}
	return f
}

// Type gets the type of the Receiver config created by this factory.
func (f *factory) Type() config.Type {
	return f.cfgType
}

// CreateDefaultConfig creates the default configuration for receiver.
func (f *factory) CreateDefaultConfig() config.Receiver {
	return f.createDefaultConfig()
}

// CreateTraceReceiver creates a component.TracesReceiver based on this config.
func (f *factory) CreateTracesReceiver(
	ctx context.Context,
	params component.ReceiverCreateParams,
	cfg config.Receiver,
	nextConsumer consumer.Traces) (component.TracesReceiver, error) {
	if f.createTraceReceiver != nil {
		return f.createTraceReceiver(ctx, params, cfg, nextConsumer)
	}
	return nil, componenterror.ErrDataTypeIsNotSupported
}

// CreateMetricsReceiver creates a component.MetricsReceiver based on this config.
func (f *factory) CreateMetricsReceiver(
	ctx context.Context,
	params component.ReceiverCreateParams,
	cfg config.Receiver,
	nextConsumer consumer.Metrics) (component.MetricsReceiver, error) {
	if f.createMetricsReceiver != nil {
		return f.createMetricsReceiver(ctx, params, cfg, nextConsumer)
	}
	return nil, componenterror.ErrDataTypeIsNotSupported
}

// CreateLogsReceiver creates a metrics processor based on this config.
func (f *factory) CreateLogsReceiver(
	ctx context.Context,
	params component.ReceiverCreateParams,
	cfg config.Receiver,
	nextConsumer consumer.Logs,
) (component.LogsReceiver, error) {
	if f.createLogsReceiver != nil {
		return f.createLogsReceiver(ctx, params, cfg, nextConsumer)
	}
	return nil, componenterror.ErrDataTypeIsNotSupported
}
