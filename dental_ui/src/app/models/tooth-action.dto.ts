import * as Immutable from 'immutable';

export interface IToothActionData {
    id: string,
    patientId: string,
    toothNo: string,
    actionId: number,
    date: Date
}

const ToothActionRecord = Immutable.Record({
    id: undefined,
    patientId: undefined,
    toothNo: undefined,
    actionId: undefined,
    date: undefined
});

export class ToothActionDto extends ToothActionRecord implements IToothActionData {
    private static fEmpty: ToothActionDto;
    
    id: string;
    patientId: string;
    toothNo: string;
    actionId: number;
    date: Date;

    constructor(data: IToothActionData) {
        super({
            id: data.id,
            patientId: data.patientId,
            toothNo: data.toothNo,
            actionId: data.actionId,
            date: data.date   
        });
    }

    static create(
        id: string,
        patientId: string,
        toothNo: string,
        actionId: number,
        date: Date): ToothActionDto {
        return new ToothActionDto({
            id,
            patientId,
            toothNo,
            actionId,
            date});
    }

    static createFromObject(obj: any): ToothActionDto {
        let data = (typeof obj === 'string') ? JSON.parse(obj) : obj;
        return new ToothActionDto({
            id: data.id !== undefined ? data.id : data.id,
            patientId: data.patientId !== undefined ? data.patientId : data.patientId,
            toothNo: data.toothNo !== undefined ? data.toothNo : data.toothNo,
            actionId: data.actionId !== undefined ? data.actionId : data.actionId,
            date: data.date !== undefined ? data.date : data.date,
        });
    }

    static get empty(): ToothActionDto {
        return ToothActionDto.fEmpty || (ToothActionDto.fEmpty = new ToothActionDto({
            id: undefined,
            patientId: undefined,
            toothNo: undefined,
            actionId: undefined,
            date: undefined
        }));
    }
}
