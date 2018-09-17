import * as Immutable from 'immutable';
import { IToothActionData } from './tooth-action.dto';

export interface ITeethData {
    diagnosisList: IToothActionData[],
    manipulationList: IToothActionData[],
}