import { ManipulationDto } from './manipulation.dto';
import { KeyedModel } from './model';
import { IActionTypeModel } from './action-type.model'

export class ManipulationModel extends KeyedModel<ManipulationDto> implements IActionTypeModel {
    private static fEmpty: ManipulationModel;

    get id(): number {
        return this.dto.id;
    }

    get name(): string {
        return this.dto.name;
    }

    get changeStatus(): number {
        return this.dto.changeStatus;
    }

    constructor(protected dto: ManipulationDto) {
        super(dto, []);
    }

    static get empty(): ManipulationModel {
        return ManipulationModel.fEmpty || (ManipulationModel.fEmpty = new ManipulationModel(ManipulationDto.createFromObject({})));
    }
}
