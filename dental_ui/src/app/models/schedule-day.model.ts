import { KeyedModel } from './model';
import { ScheduleByDayDto, IAppointmentData } from './schedule-day.dto';

export class ScheduleByDayModel extends KeyedModel<ScheduleByDayDto> {
    private static fEmpty: ScheduleByDayModel;

    get day(): Date {
        return this.dto.day;
    }

    get appointments(): IAppointmentData[] {
        return this.dto.appointments;
    }

    constructor(protected dto: ScheduleByDayDto) {
        super(dto, []);
    }

    static get empty(): ScheduleByDayModel {
        return ScheduleByDayModel.fEmpty || (ScheduleByDayModel.fEmpty = new ScheduleByDayModel(ScheduleByDayDto.createFromObject({})));
    }
}
