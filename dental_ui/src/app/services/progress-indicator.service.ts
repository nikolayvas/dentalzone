import { Injectable } from '@angular/core';

import { ProgressIndicatorActions } from '../store/progress-indicator.actions';
import { StoreService } from '../services/store.service';

@Injectable()
export class ProgressIndicatorService {

    constructor(private storeService: StoreService) {
    }

    set(state: { isActive?: boolean, freezeUI?: boolean, progress?: number, infinite?: boolean }) {
        this.storeService.dispatch({
            type: ProgressIndicatorActions.SET_PROGRESS_INDICATOR,
            payload: { isActive: state.isActive, freezeUI: state.freezeUI, progress: state.progress, infinite: state.infinite }
        });
    }
}