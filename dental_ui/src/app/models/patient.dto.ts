import * as Immutable from 'immutable';

export interface IPatientData {
    id: string,
    firstName: string,
    middleName: string,
    lastName: string,
    email: string,
    address: string,
    phoneNumber: string,
    generalInfo: string,
    registrationDate: Date
}

const CompanyRecord = Immutable.Record({
    id: undefined,
    firstName: undefined,
    middleName: undefined,
    lastName: undefined,
    email: undefined,
    address: undefined,
    phoneNumber: undefined,
    generalInfo: undefined,
    registrationDate: undefined
});

export class PatientDto extends CompanyRecord implements IPatientData {
    private static fEmpty: PatientDto;
    
    id: string;
    firstName: string;
    middleName: string;
    lastName: string;
    email: string;
    address: string;
    phoneNumber: string;
    generalInfo: string;
    registrationDate: Date;

    constructor(data: IPatientData) {
        super({
            id: data.id,
            firstName: data.firstName,
            middleName: data.middleName,
            lastName: data.lastName,
            email: data.email,
            address: data.address,
            phoneNumber: data.phoneNumber,
            generalInfo: data.generalInfo,
            registrationDate: data.registrationDate
        });
    }

    static create(
        id: string,
        firstName: string,
        middleName: string,
        lastName: string,
        email: string,
        address: string,
        phoneNumber: string,
        generalInfo: string,
        registrationDate: Date): PatientDto {
        return new PatientDto({
            id,
            firstName,
            middleName,
            lastName,
            email,
            address,
            phoneNumber,
            generalInfo,
            registrationDate
        });
    }

    static createFromObject(obj: any): PatientDto {
        let data = (typeof obj === 'string') ? JSON.parse(obj) : obj;
        return new PatientDto({
            id: data.id !== undefined ? data.id : data.Id,
            firstName: data.firstName !== undefined ? data.firstName : data.FirstName,
            middleName: data.middleName !== undefined ? data.middleName : data.MiddleName,
            lastName: data.lastName !== undefined ? data.lastName : data.LastName,
            email: data.email !== undefined ? data.email : data.email,
            address: data.address !== undefined ? data.address : data.Address,
            phoneNumber: data.phoneNumber !== undefined ? data.phoneNumber : data.PhoneNumber,
            generalInfo: data.generalInfo !== undefined ? data.generalInfo : data.GeneralInfo,
            registrationDate: data.registrationDate !== undefined ? data.registrationDate : data.RegistrationDate
        });
    }

    static get empty(): PatientDto {
        return PatientDto.fEmpty || (PatientDto.fEmpty = new PatientDto({
            id: undefined,
            firstName: undefined,
            middleName: undefined,
            lastName: undefined,
            email: undefined,
            address: undefined,
            phoneNumber: undefined,
            generalInfo: undefined,
            registrationDate: undefined
        }));
    }
}
