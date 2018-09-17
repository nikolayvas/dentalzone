import { PatientDto } from './patient.dto';
import { KeyedModel } from './model';

export class PatientModel extends KeyedModel<PatientDto> {
    private static fEmpty: PatientModel;

    get id(): string {
        return this.dto.id;
    }

    get firstName(): string {
        return this.dto.firstName;
    }

    get middleName(): string {
        return this.dto.middleName;
    }

    get lastName(): string {
        return this.dto.lastName;
    }

    get email(): string {
        return this.dto.email;
    }

    get address(): string {
        return this.dto.address;
    }

    get phoneNumber(): string {
        return this.dto.phoneNumber;
    }

    get generalInfo(): string {
        return this.dto.generalInfo;
    }

    get registrationDate(): Date {
        return this.dto.registrationDate;
    }

    constructor(protected dto: PatientDto) {
        super(dto, []);
    }

    static get empty(): PatientModel {
        return PatientModel.fEmpty || (PatientModel.fEmpty = new PatientModel(PatientDto.createFromObject({})));
    }
}
