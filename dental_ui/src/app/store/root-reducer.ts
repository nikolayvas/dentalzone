import { combineReducers, ActionReducerMap } from '@ngrx/store';

import { ClientPortalSeedActions } from '../store/seed-actions';
import { patientReducer, IPatientsStoreState } from './patient-reducer';
import { metaDataReducer, IMetaDataState } from './seed-reducer';
import { teethDataReducer, ITeethDataState } from './teeth-data.reducer';

/**
 * Client Portal combined state.
 */
export interface IClientPortalStoreState {
    patientsState: IPatientsStoreState,
    metaData: IMetaDataState,
    teethDataState: ITeethDataState
};

/**
 * Client Portal combined reducer.
 */
export const clientPortalStoreReducer = combineReducers({
    patientsState: patientReducer,
    metaData: metaDataReducer,
    teethDataState: teethDataReducer
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
