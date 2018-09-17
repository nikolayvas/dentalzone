import { StateManager } from './state-manager';
import * as R from 'ramda';

export class StateManagerHelper {
    static updateCollection<TData, TDto, TModel extends { id: string }>(
        models: TModel[],
        payload: any,
        stateManager: StateManager<TData, TDto, TModel>): TModel[] {

        const existingModel = R.find(x => x.id === payload.id, models);
        const existingModelIndex = models.indexOf(existingModel);
        const updatedModel = stateManager.updateModel(payload, existingModel);

        return R.update(existingModelIndex, updatedModel, models);
    }
}

