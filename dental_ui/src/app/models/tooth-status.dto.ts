import * as Immutable from 'immutable';

export interface IToothStatusData {
    id: string,
    name: string,
}

const ToothStatusRecord = Immutable.Record({
    id: undefined,
    name: undefined,
});

export class ToothStatusDto extends ToothStatusRecord implements IToothStatusData {
    private static fEmpty: ToothStatusDto;
    
    id: string;
    name: string;

    constructor(data: IToothStatusData) {
        super({
            id: data.id,
            name: data.name,
        });
    }

    static create(
        id: string,
        name: string,
        changeStatus: number): ToothStatusDto {
        return new ToothStatusDto({
            id,
            name});
    }

    static createFromObject(obj: any): ToothStatusDto {
        let data = (typeof obj === 'string') ? JSON.parse(obj) : obj;
        return new ToothStatusDto({
            id: data.id !== undefined ? data.id : data.id,
            name: data.name !== undefined ? data.name : data.name,
        });
    }

    static get empty(): ToothStatusDto {
        return ToothStatusDto.fEmpty || (ToothStatusDto.fEmpty = new ToothStatusDto({
            id: undefined,
            name: undefined,
        }));
    }
}
