import { ActionReducer } from '@ngrx/store';

import * as R from 'ramda';

import { PatientActions } from './patient-actions';
import { IPatientData, PatientDto } from '../models/patient.dto';
import { PatientModel } from '../models/patient.model';

import { PayloadAction } from './payload-action';
import { StateManager } from './state-manager';
import { StateManagerHelper } from './state-manager-helper';

/**
 * Patients state interface.
 */
export interface IPatientsStoreState {
    patients: PatientModel[];
};

/**
 * Patients reducer.
 */
export const patientReducer: ActionReducer<IPatientsStoreState> = (state: IPatientsStoreState = { patients: undefined }, action: PayloadAction) => {

    const stateManager = new StateManager<IPatientData, PatientDto, PatientModel>(PatientDto, PatientModel);

    switch (action.type) {
        case PatientActions.SEED_PATIENTS:
            const seedModel = stateManager.createModels(action.payload);
            return { patients: seedModel };

        case PatientActions.CREATE_PATIENT:
            const patient = stateManager.createModel(action.payload);
            return { patients: [...state.patients, patient] };

        case PatientActions.UPDATE_PATIENT:
            const updatedModel = StateManagerHelper.updateCollection(state.patients, action.payload, stateManager);
            return { patients: updatedModel };
        
        case PatientActions.REMOVE_PATIENT:
        const toRemove = state.patients.find(x => x.id === action.payload);
        return {
            patients: [...R.without([toRemove], state.patients)]
        };

        default:
            return state;
    }
}