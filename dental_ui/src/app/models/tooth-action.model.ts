import { ToothActionDto } from './tooth-action.dto';
import { KeyedModel } from './model';

export class ToothActionModel extends KeyedModel<ToothActionDto> {
    private static fEmpty: ToothActionModel;

    get id(): string {
        return this.dto.id;
    }

    get patientId(): string {
        return this.dto.patientId;
    }
    get toothNo(): string {
        return this.dto.toothNo;
    }

    get actionId(): number {
        return this.dto.actionId;
    }

    get date(): Date {
        return this.dto.date;
    }

    constructor(protected dto: ToothActionDto) {
        super(dto, []);
    }

    static get empty(): ToothActionModel {
        return ToothActionModel.fEmpty || (ToothActionModel.fEmpty = new ToothActionModel(ToothActionDto.createFromObject({})));
    }
}
