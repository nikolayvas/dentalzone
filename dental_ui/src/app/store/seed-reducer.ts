import { ActionReducer, Action } from '@ngrx/store';
import { ManipulationModel } from '../models/manipulation.model';
import { ToothStatusModel } from '../models/tooth-status.model';
import { DiagnosisModel } from '../models/diagnosis.mode';
import { PayloadAction } from './payload-action';
import { StateManager } from './state-manager';

import { IToothStatusData, ToothStatusDto } from '../models/tooth-status.dto'
import { IManipulationData, ManipulationDto } from '../models/manipulation.dto'
import { IDiagnosisData, DiagnosisDto } from '../models/diagnosis.dto'
import { ClientPortalSeedActions } from './seed-actions'
import { IMetaData } from '../models/meta-data.dto'
/**
 * Enumerations state interface.
 */
export interface IMetaDataState {
    toothStatuses: ToothStatusModel[];
    manipupations: ManipulationModel[];
    diagnosis: DiagnosisModel[];
};

/**
 * Enumerations reducer.
 */
export const metaDataReducer: ActionReducer<IMetaDataState> = (state: IMetaDataState = 
    {toothStatuses: [], manipupations: [], diagnosis: []}, action: PayloadAction) => {

    const stateToothStatuses = new StateManager<IToothStatusData, ToothStatusDto, ToothStatusModel>(ToothStatusDto, ToothStatusModel);
    const stateManipupations = new StateManager<IManipulationData, ManipulationDto, ManipulationModel>(ManipulationDto, ManipulationModel);
    const stateDiagnosis = new StateManager<IDiagnosisData, DiagnosisDto, DiagnosisModel>(DiagnosisDto, DiagnosisModel);

    switch (action.type) {
        case ClientPortalSeedActions.RESET_STATE_DATA:

        const data = action.payload as IMetaData;
        const seedState = {
            toothStatuses: data.toothStatusList.map(x => stateToothStatuses.createModel(x)),
            manipupations: data.manipulationList.map(x => stateManipupations.createModel(x)),
            diagnosis: data.diagnosisList.map(x => stateDiagnosis.createModel(x))
        };
        
        return seedState;

        default:
            return state;
    }
}