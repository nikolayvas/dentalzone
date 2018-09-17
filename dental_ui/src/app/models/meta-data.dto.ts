import * as Immutable from 'immutable';
import { IDiagnosisData } from './diagnosis.dto';
import { IManipulationData } from './manipulation.dto';
import { IToothStatusData } from './tooth-status.dto';

export interface IMetaData {
    diagnosisList: IDiagnosisData[],
    manipulationList: IManipulationData[],
    toothStatusList: IToothStatusData[]
}
