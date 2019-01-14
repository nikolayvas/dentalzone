import { combineReducers, ActionReducerMap } from '@ngrx/store';

import { patientReducer, IPatientsStoreState } from './patient-reducer';
import { metaDataReducer, IMetaDataState } from './seed-reducer';
import { teethDataReducer, ITeethDataState } from './teeth-data.reducer';
import { IUIStoreState, uiReducer } from './ui-reducer';
import { progressIndicatorReducer, IProgressIndicatorState } from './progress-indicator.reducer';

/**
 * Client Portal combined state.
 */
export interface IClientPortalStoreState {
    uiState: IUIStoreState,
    patientsState: IPatientsStoreState,
    metaData: IMetaDataState,
    teethDataState: ITeethDataState,
    progressIndicatorState: IProgressIndicatorState,
};

/**
 * Client Portal combined reducer.
 */
export const clientPortalStoreReducer = combineReducers({
    uiState: uiReducer,
    patientsState: patientReducer,
    metaData: metaDataReducer,
    teethDataState: teethDataReducer,
    progressIndicatorState: progressIndicatorReducer,
});

/**
 * Client portal data state.
 */
export interface IClientPortalDataState {
    clientPortalStore: IClientPortalStoreState
}

/**
 * Reducers for the IClientPortalDataState.
 */
const clientPortalDataState = combineReducers({
    clientPortalStore: clientPortalStoreReducer
});

export function clientPortalDataReducer(state, action) {
    return clientPortalDataState(state, action);
};

export class ClientPortalState {
    data: IClientPortalDataState;
}

export const rootReducer: ActionReducerMap<ClientPortalState> = {
    data: clientPortalDataReducer,
}