import { Iterable } from 'immutable';
import * as R from 'ramda';

//
// A helper class for:
//  1. Creating models from given key/value pairs data.
//  2. Updating existing (given) model with values from given key/value pairs.
//     The model is actually not updated/mutated, but a new one is returned, as the flux pattern suggests.
//  3. Seeding models.
// See loan-details.reducers.ts for an example usage.
//
export class StateManager<TData, TDto, TModel> {

    constructor(
        private dtoConstructor: { new (data: TData): TDto },
        private modelConstructor: { new (dto: TDto): TModel }) {
    }

    createModel(data: TData): TModel {

        return data ? new this.modelConstructor(new this.dtoConstructor(data)) : undefined;
    }

    createModels(data: TData[]): TModel[] {

        return data ? data.map(x => this.createModel(x)) : undefined;
    }

    updateModel(data: any, existingModel: TModel): TModel {
        const keys = Object.keys(existingModel);
        const existingData: { [key: string]: any } = {};
        keys.forEach(key => existingData[key] = existingModel[key]);
        const newData = <TData>Object.assign({}, existingData, data);
        return this.createModel(newData);
    }
}