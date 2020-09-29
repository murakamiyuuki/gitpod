// Copyright (c) 2020 TypeFox GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package supervisor

import (
	"context"
	"os"
	"os/exec"
	"sync"

	"github.com/gitpod-io/gitpod/common-go/log"
	"github.com/gitpod-io/gitpod/supervisor/api"
)

type tasksManager struct {
	cfg   *Config
	tasks map[string]*api.TasksStatus
	Ready chan struct{}
}

func newTasksManager(cfg *Config) *tasksManager {
	return &tasksManager{
		cfg:   cfg,
		tasks: make(map[string]*api.TasksStatus),
		Ready: make(chan struct{}),
	}
}

func (tm *tasksManager) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(tm.Ready)

	supervisorExecutable, err := os.Executable()
	if err != nil {
		log.WithError(err).Warn("cannot find supervisor executable")
		return
	}

	newTerminalOutput, err := exec.Command(supervisorExecutable, "terminal", "new", "-d").CombinedOutput()
	if err != nil {
		log.WithError(err).Warn("cannot craete a new terminal")
		return
	}
	alias := string(newTerminalOutput)
	tm.tasks[alias] = &api.TasksStatus{
		Alias: alias,
	}
}
