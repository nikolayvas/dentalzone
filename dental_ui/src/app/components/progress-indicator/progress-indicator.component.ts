import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs';

import { IProgressIndicatorState } from '../../store/progress-indicator.reducer';
import { Utils } from '../../services/utils';
import { StoreService } from '../../services/store.service';


@Component({
    selector: 'progress-indicator',
    templateUrl: 'progress-indicator.component.html',
    styleUrls: ['progress-indicator.component.css']
})
export class ProgressIndicatorComponent implements OnInit, OnDestroy {

    private _subscriptions: Subscription[] = [];
    private _state: IProgressIndicatorState;

    public get isActive(): boolean {
        return this._state && this._state.isActive;
    }

    public get freezeUI(): boolean {
        return this._state && this._state.freezeUI;
    }

    public get progress(): number {
        return this._state ? Math.round(this._state.progress) : 0;
    }

    public get infinite(): boolean {
        return this._state && this._state.infinite;
    }

    constructor(private store: StoreService) {
    }

    ngOnInit() {
        this._subscriptions.push(
            this.store.select(x => x.data.clientPortalStore.progressIndicatorState).subscribe(state => this._state = state)
        );
    }

    ngOnDestroy() {
        Utils.unsubscribe(this._subscriptions);
    }
}
