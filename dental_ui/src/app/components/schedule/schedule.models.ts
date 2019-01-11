import * as moment from 'moment';

export enum DayOrWeekMode {
    Day = 1,
    Week = 2,
};

export enum WeekDay{
    Monday = 1,
    Tuesday = 2,
    Wednesday = 3,
    Thursday = 4,
    Friday = 5,
    Saturday = 6,
    Sunday = 7
};

export interface IScheduleRowModel {
    //day
    //hour
    time: string,
};

export interface IDayOfWeekModel {
    dayOfWeek: string,
    dayOfMonth: string,
    isToday: boolean
};

export interface IPaginatorModel {
    pageMode: DayOrWeekMode,
    currentDate: moment.Moment 
};

export interface IAppointmentModel {
    x: number,
    y: number,
    dateTime: Date,
    hasNext: boolean,
    patientID: string | null,
    patientName: string | null,
    color: string
};
