import { ActionReducer } from '@ngrx/store';

import { UIActions } from './ui-actions';

import { PayloadAction } from './payload-action';

/**
 * UI state interface.
 */
export interface IUIStoreState {
    patients_filter: string;
};

/**
 * UI reducer.
 */
export const uiReducer: ActionReducer<IUIStoreState> = (state: IUIStoreState = { patients_filter: undefined }, action: PayloadAction) => {

    switch (action.type) {
        case UIActions.SET_PATIENTS_FILTER:
            return { patients_filter: action.payload };

        default:
            return state;
    }
}