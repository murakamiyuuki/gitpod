/**
 * Copyright (c) 2020 TypeFox GmbH. All rights reserved.
 * Licensed under the GNU Affero General Public License (AGPL).
 * See License-AGPL.txt in the project root for license information.
 */

import { injectable, inject } from 'inversify';
import { FrontendApplicationContribution } from '@theia/core/lib/browser';
import type { TasksStatusResponse } from '@gitpod/supervisor-api-grpc/lib/status_pb';
import { TerminalFrontendContribution } from '@theia/terminal/lib/browser/terminal-frontend-contribution';

@injectable()
export class GitpodTaskContribution implements FrontendApplicationContribution {

    @inject(TerminalFrontendContribution)
    protected readonly terminals: TerminalFrontendContribution;

    initializeLayout() {
        fetch(window.location.protocol + '//' + window.location.host + '/_supervisor/v1/status/tasks').then(async response => {
            const status: TasksStatusResponse.AsObject = await response.json();
            for (const task of status.tasksList) {
                const terminal = await this.terminals.newTerminal({
                    id: 'gitpod-task:' + task.alias
                });
                await terminal.start();
                terminal.sendText(`/theia/supervisor terminal attach ${task.alias} -ir`);
                this.terminals.open(terminal, { mode: 'open' });
            }
        }).catch(err =>
            console.error('Failed to initialize Gitpod tasks:', err)
        );
    }

}