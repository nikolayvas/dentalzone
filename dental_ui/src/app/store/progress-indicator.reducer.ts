import { ActionReducer } from '@ngrx/store';
import { PayloadAction } from './payload-action';
import { ProgressIndicatorActions } from './progress-indicator.actions';

export interface IProgressIndicatorState {
    isActive: boolean;
    freezeUI: boolean;
    progress: number;
    infinite: boolean;
}

export const progressIndicatorReducer: ActionReducer<IProgressIndicatorState> =
    (state: IProgressIndicatorState = { isActive: false, freezeUI: false, progress: 0, infinite: false }, action: PayloadAction) => {

        switch (action.type) {
            case ProgressIndicatorActions.SET_PROGRESS_INDICATOR:
                return {
                    isActive: action.payload.isActive !== undefined ? action.payload.isActive : state.isActive,
                    freezeUI: action.payload.freezeUI !== undefined ? action.payload.freezeUI : state.freezeUI,
                    progress: action.payload.progress !== undefined ? action.payload.progress : state.progress,
                    infinite: action.payload.infinite !== undefined ? action.payload.infinite : state.infinite,
                }

            default:
                return state;
        }
    };

