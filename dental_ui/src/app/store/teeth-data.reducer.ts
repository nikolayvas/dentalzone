import { ActionReducer, Action } from '@ngrx/store';
import { ToothActionModel } from '../models/tooth-action.model';
import { PayloadAction } from './payload-action';
import { StateManager } from './state-manager';

import { IToothActionData, ToothActionDto } from '../models/tooth-action.dto'
import { ITeethData } from '../models/teeth-data.dto'
import { TeethDataActions } from './teeth-data.action'
/**
 * Enumerations state interface.
 */
export interface ITeethDataState {
    toothDiagnosis: ToothActionModel[];
    toothManipupations : ToothActionModel[];
};

/**
 * Teeth data reducer.
 */
export const teethDataReducer: ActionReducer<ITeethDataState> = (state: ITeethDataState = 
    {toothDiagnosis: [], toothManipupations: []}, action: PayloadAction) => {

    const stateToothDiagnosis = new StateManager<IToothActionData, ToothActionDto, ToothActionModel>(ToothActionDto, ToothActionModel);
    const stateToothManipulation = new StateManager<IToothActionData, ToothActionDto, ToothActionModel>(ToothActionDto, ToothActionModel);

    switch (action.type) {
        case TeethDataActions.SEED_TEETH_STATE_DATA:

            const data = action.payload as ITeethData;
            const seedState = {
                toothDiagnosis: data.diagnosisList.map(x => stateToothDiagnosis.createModel(x)),
                toothManipupations: data.manipulationList.map(x => stateToothManipulation.createModel(x)),
            };
        
            return seedState;
        case TeethDataActions.CLEAR_TEETH_STATE_DATA:
            return  {
                toothDiagnosis: [],
                toothManipupations: [],
            };

        case TeethDataActions.ADD_DIAGNOSIS:

            const newDiagnosis = stateToothDiagnosis.createModel(action.payload);
            
            return { 
                toothDiagnosis : [...state.toothDiagnosis, newDiagnosis],
                toothManipupations: state.toothManipupations
            };

        case TeethDataActions.ADD_MANIPULATION:

            const newManipulation = stateToothManipulation.createModel(action.payload);
            
            return { 
                toothDiagnosis : state.toothDiagnosis,
                toothManipupations: [...state.toothManipupations, newManipulation]
            };

        case TeethDataActions.REMOVE_DIAGNOSIS:

            const toothDiagnosisID = action.payload;
            
            return { 
                toothDiagnosis : [...state.toothDiagnosis.filter(n=>n.id !== toothDiagnosisID)],
                toothManipupations: state.toothManipupations
            };

        case TeethDataActions.REMOVE_MANIPULATION:

            const toothManipulationID = action.payload;
            
            return { 
                toothDiagnosis : state.toothDiagnosis,
                toothManipupations: [...state.toothManipupations.filter(n=>n.id !== toothManipulationID)]
            };

        default:
            return state;
    }
}