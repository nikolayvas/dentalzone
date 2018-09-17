import * as Immutable from 'immutable';

export interface IDiagnosisData {
    id: number,
    name: string,
    changeStatus: number
}

const DiagnosisRecord = Immutable.Record({
    id: undefined,
    name: undefined,
    changeStatus: undefined
});

export class DiagnosisDto extends DiagnosisRecord implements IDiagnosisData {
    private static fEmpty: DiagnosisDto;
    
    id: number;
    name: string;
    changeStatus: number;

    constructor(data: IDiagnosisData) {
        super({
            id: data.id,
            name: data.name,
            changeStatus: data.changeStatus   
        });
    }

    static create(
        id: number,
        name: string,
        changeStatus: number): DiagnosisDto {
        return new DiagnosisDto({
            id,
            name,
            changeStatus});
    }

    static createFromObject(obj: any): DiagnosisDto {
        let data = (typeof obj === 'string') ? JSON.parse(obj) : obj;
        return new DiagnosisDto({
            id: data.id !== undefined ? data.id : data.id,
            name: data.name !== undefined ? data.name : data.name,
            changeStatus: data.changeStatus !== undefined ? data.changeStatus : data.changeStatus,
        });
    }

    static get empty(): DiagnosisDto {
        return DiagnosisDto.fEmpty || (DiagnosisDto.fEmpty = new DiagnosisDto({
            id: undefined,
            name: undefined,
            changeStatus: undefined
        }));
    }
}
