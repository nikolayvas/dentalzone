import * as Immutable from 'immutable';

export interface IManipulationData {
    id: number,
    name: string,
    changeStatus: number
}

const ManipulationRecord = Immutable.Record({
    id: undefined,
    name: undefined,
    changeStatus: undefined
});

export class ManipulationDto extends ManipulationRecord implements IManipulationData {
    private static fEmpty: ManipulationDto;
    
    id: number;
    name: string;
    changeStatus: number;

    constructor(data: IManipulationData) {
        super({
            id: data.id,
            name: data.name,
            changeStatus: data.changeStatus   
        });
    }

    static create(
        id: number,
        name: string,
        changeStatus: number): ManipulationDto {
        return new ManipulationDto({
            id,
            name,
            changeStatus});
    }

    static createFromObject(obj: any): ManipulationDto {
        let data = (typeof obj === 'string') ? JSON.parse(obj) : obj;
        return new ManipulationDto({
            id: data.id !== undefined ? data.id : data.id,
            name: data.name !== undefined ? data.name : data.name,
            changeStatus: data.changeStatus !== undefined ? data.changeStatus : data.changeStatus,
        });
    }

    static get empty(): ManipulationDto {
        return ManipulationDto.fEmpty || (ManipulationDto.fEmpty = new ManipulationDto({
            id: undefined,
            name: undefined,
            changeStatus: undefined
        }));
    }
}
