import * as moment from 'moment';

export enum DayOrWeekMode {
    Day = 1,
    Week = 2,
};

export interface IScheduleRowModel {
    //day
    //hour
    time: string,
    forDays: IAppointmentModel[]
};

export interface IDayOfWeekModel {
    dayOfWeek: string,
    dayOfMonth: string,
    isToday: boolean
};

export interface IPaginatorModel {
    pageMode: number,
    currentDate: moment.Moment 
};

export interface IAppointmentModel {
    x: number,
    y: number,
    hasPrev: boolean,
    hasNext: boolean,
    patientID: string | null,
    patientName: string | null,
    color: string
};
