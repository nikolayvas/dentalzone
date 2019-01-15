

import { ActionReducer } from '@ngrx/store';
import { PayloadAction } from './payload-action';
import { ScheduleActions } from './schedule.actions';
import { StateManager } from './state-manager';
import { ScheduleByDayModel } from '../models/schedule-day.model';
import { IScheduleByDayData, ScheduleByDayDto } from '../models/schedule-day.dto';

import * as R from 'ramda';

export interface IScheduleState {
    perDay: ScheduleByDayModel[];
}

export const scheduleReducer: ActionReducer<IScheduleState> =
    (state: IScheduleState = { perDay: []}, action: PayloadAction) => {
        const stateManager = new StateManager<IScheduleByDayData, ScheduleByDayDto, ScheduleByDayModel>(ScheduleByDayDto, ScheduleByDayModel);
        
        switch (action.type) {
            case ScheduleActions.SEED_DATA_FOR_DAY:

                const seedModel = stateManager.createModel(action.payload);

                return { perDay: [...state.perDay, seedModel] };
                
            case ScheduleActions.UPDATE_DATA_FOR_DAY: 
                const updateModel = stateManager.createModel(action.payload);

                const existingModelIndex = state.perDay.findIndex(n=>
                    n.day.getFullYear() === updateModel.day.getFullYear() && 
                    n.day.getMonth() === updateModel.day.getMonth() &&
                    n.day.getDate() === updateModel.day.getDate()
                    );

                return { perDay: R.update(existingModelIndex, updateModel, state.perDay) };


            default:
                return state;
        }
    };

