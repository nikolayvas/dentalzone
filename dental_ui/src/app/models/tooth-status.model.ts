import { ToothStatusDto } from './tooth-status.dto';
import { KeyedModel } from './model';

export class ToothStatusModel extends KeyedModel<ToothStatusDto> {
    private static fEmpty: ToothStatusModel;

    get id(): string {
        return this.dto.id;
    }

    get name(): string {
        return this.dto.name;
    }

    constructor(protected dto: ToothStatusDto) {
        super(dto, []);
    }

    static get empty(): ToothStatusModel {
        return ToothStatusModel.fEmpty || (ToothStatusModel.fEmpty = new ToothStatusModel(ToothStatusDto.createFromObject({})));
    }
}
