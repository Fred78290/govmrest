/*
Copyright (c) 2023 Fred78290, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Based on govc sources https://github.com/vmware/govmomi/govc
*/

package flags

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Fred78290/govmrest/vim25"
	"github.com/Fred78290/vmrest-go-client/client"
	"github.com/vmware/govmomi/govc/flags"
)

const (
	envURL      = "GOVMREST_URL"
	envUsername = "GOVMREST_USERNAME"
	envPassword = "GOVMREST_PASSWORD"
	envTimeout  = "GOVMREST_TIMEOUT"
)

const cDescr = "vmrest URL"

type ClientFlag struct {
	common

	*flags.DebugFlag

	endpoint *url.URL
	username string
	password string
	timeout  time.Duration
	client   *vim25.Client
}

var (
	home          = os.Getenv("GOVMREST_HOME")
	clientFlagKey = flagKey("client")
)

func init() {
	if home == "" {
		home = filepath.Join(os.Getenv("HOME"), ".govmrest")
	}
}

func NewClientFlag(ctx context.Context) (*ClientFlag, context.Context) {
	if v := ctx.Value(clientFlagKey); v != nil {
		return v.(*ClientFlag), ctx
	}

	u, _ := url.Parse("https://localhost:8697")
	v := &ClientFlag{
		endpoint: u,
		timeout:  120,
	}
	v.DebugFlag, ctx = flags.NewDebugFlag(ctx)
	ctx = context.WithValue(ctx, clientFlagKey, v)
	return v, ctx
}

func (flag *ClientFlag) String() string {
	return flag.endpoint.String()
}

func (flag *ClientFlag) Set(s string) error {
	var err error

	flag.endpoint, err = url.Parse(s)

	return err
}

func (flag *ClientFlag) Register(ctx context.Context, f *flag.FlagSet) {
	flag.RegisterOnce(func() {
		flag.DebugFlag.Register(ctx, f)

		{
			flag.Set(os.Getenv(envURL))
			usage := fmt.Sprintf("%s [%s]", cDescr, envURL)
			f.Var(flag, "u", usage)
		}

		{
			flag.username = os.Getenv(envUsername)
			flag.password = os.Getenv(envPassword)
			if timeout, err := time.ParseDuration(os.Getenv(envTimeout)); err == nil {
				flag.timeout = timeout
			}
		}
	})
}

func (flag *ClientFlag) Process(ctx context.Context) error {
	return flag.ProcessOnce(func() error {
		err := flag.DebugFlag.Process(ctx)
		if err != nil {
			return err
		}

		// Override username if set
		if flag.username != "" {
			var password string
			var ok bool

			if flag.endpoint.User != nil {
				password, ok = flag.endpoint.User.Password()
			}

			if ok {
				flag.endpoint.User = url.UserPassword(flag.username, password)
			} else {
				flag.endpoint.User = url.User(flag.username)
			}
		}

		// Override password if set
		if flag.password != "" {
			var username string

			if flag.endpoint.User != nil {
				username = flag.endpoint.User.Username()
			}

			flag.endpoint.User = url.UserPassword(username, flag.password)
		}

		return nil
	})
}

func (flag *ClientFlag) Client() (*vim25.Client, error) {
	if flag.client != nil {
		return flag.client, nil
	}

	cfg := &client.Configuration{
		Endpoint:  flag.endpoint.String(),
		UserAgent: "govmrest/1.0.0/go",
		UserName:  flag.username,
		Password:  flag.password,
		Timeout:   120,
	}

	if client, err := client.NewAPIClient(cfg); err != nil {
		return nil, err
	} else {
		flag.client = &vim25.Client{
			APIClient: client,
		}
	}

	return flag.client, nil
}

// Environ returns the govc environment variables for this connection
func (flag *ClientFlag) Environ(extra bool) []string {
	var env []string
	add := func(k, v string) {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	u := *flag.endpoint
	if u.User != nil {
		add(envUsername, u.User.Username())

		if p, ok := u.User.Password(); ok {
			add(envPassword, p)
		}

		u.User = nil
	}

	u.Fragment = ""
	u.RawQuery = ""

	add(envURL, strings.TrimPrefix(u.String(), "https://"))

	return env
}

// WithCancel calls the given function, returning when complete or canceled via SIGINT.
func (flag *ClientFlag) WithCancel(ctx context.Context, f func(context.Context) error) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)

	wctx, cancel := context.WithCancel(ctx)
	defer cancel()

	done := make(chan bool)
	var werr error

	go func() {
		defer close(done)
		werr = f(wctx)
	}()

	select {
	case <-sig:
		cancel()
		<-done // Wait for f() to complete
	case <-done:
	}

	return werr
}
