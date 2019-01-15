import * as Immutable from 'immutable';

export interface IAppointmentData {
    date : Date;
    patientID: string;
};

export interface IScheduleByDayData {
    day: Date,
    appointments: IAppointmentData[];
}

const ScheduleByDayRecord = Immutable.Record({
    day: undefined,
    appointments: undefined,
});

export class ScheduleByDayDto extends ScheduleByDayRecord implements IScheduleByDayData {
    private static fEmpty: ScheduleByDayDto;
    
    day: Date;
    appointments:  IAppointmentData[];

    constructor(data: IScheduleByDayData) {
        super({
            day: data.day,
            appointments: data.appointments,
        });
    }

    static create(
        day: Date,
        appointments: IAppointmentData[],
        ): ScheduleByDayDto {
        return new ScheduleByDayDto({
            day: day,
            appointments: appointments
        });
    }

    static createFromObject(obj: any): ScheduleByDayDto {
        let data = (typeof obj === 'string') ? JSON.parse(obj) : obj;
        return new ScheduleByDayDto({
            day: data.day !== undefined ? data.day : data.day,
            appointments: data.appointments !== undefined ? data.appointments : data.appointments,
        });
    }

    static get empty(): ScheduleByDayDto {
        return ScheduleByDayDto.fEmpty || (ScheduleByDayDto.fEmpty = new ScheduleByDayDto({
            day: undefined,
            appointments: [],
        }));
    }
}