import { DiagnosisDto } from './diagnosis.dto';
import { KeyedModel } from './model';
import { IActionTypeModel } from './action-type.model'

export class DiagnosisModel extends KeyedModel<DiagnosisDto> implements IActionTypeModel {
    private static fEmpty: DiagnosisModel;

    get id(): number {
        return this.dto.id;
    }

    get name(): string {
        return this.dto.name;
    }

    get changeStatus(): number {
        return this.dto.changeStatus;
    }

    constructor(protected dto: DiagnosisDto) {
        super(dto, []);
    }

    static get empty(): DiagnosisModel {
        return DiagnosisModel.fEmpty || (DiagnosisModel.fEmpty = new DiagnosisModel(DiagnosisDto.createFromObject({})));
    }
}
