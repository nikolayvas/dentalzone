import { Injectable } from '@angular/core';
import { Subscription, Observable, zip as observableZip, of as observableOf } from 'rxjs';
import { take, filter, map,  } from 'rxjs/operators';

import { StoreService } from '../../services/store.service';
import { IAppointmentModel, IPaginatorModel, DayOrWeekMode, WeekDay } from './schedule.models';

import * as moment from 'moment';

@Injectable()
export class ScheduleService {

    constructor(
        private store: StoreService
    ) {
    }

    getAppointmentsPerDay$(day: moment.Moment): Observable<IAppointmentModel[]> {

        var now = new Date(day.valueOf());
        var date = new Date(now.getFullYear(), now.getMonth(), now.getDay());
        date.setHours(date.getHours() +  6);

        var time1 = new Date(date.getTime());
        time1.setHours(time1.getHours() +  2);

        var time2 = new Date(date.getTime());
        time2.setHours(time1.getHours() +  4);
        time2.setMinutes(time1.getMinutes() +  45);

        var appointments: IAppointmentModel[] = [];
        appointments.push({
            x: undefined,
            y: undefined,
            dateTime: time1,
            color: undefined,
            hasNext: undefined,
            patientID: "5c2cce35e3bf85e693044b61",
            patientName: "ass fdfdsfdsf dfds fdsf"
        });

        appointments.push({
            x: undefined,
            y: undefined,
            dateTime: time2,
            color: undefined,
            hasNext: undefined,
            patientID: "5c2cceb4e3bf85e693044b62",
            patientName: "Nikolay"
        });

        return observableOf(appointments);
    }

    getAppointmentsPerDayOrWeek$(data: IPaginatorModel): Observable<{ [day: string] : IAppointmentModel[] }> {
        if(data.pageMode == DayOrWeekMode.Day) {
            return this.getAppointmentsPerDay$(data.currentDate.clone()).pipe(filter(n=>!!n), take(1), map(n=> {
                var appointmentsPerDay: { [day: string] : IAppointmentModel[] } = {};
                appointmentsPerDay[1] = n;

                return appointmentsPerDay;
            }));
        }
        else {
            const zip = observableZip<{ [day: string]: IAppointmentModel[] }>(
                this.getAppointmentsPerDay$(data.currentDate.clone().weekday(WeekDay.Monday)),
                this.getAppointmentsPerDay$(data.currentDate.clone().weekday(WeekDay.Tuesday)),
                this.getAppointmentsPerDay$(data.currentDate.clone().weekday(WeekDay.Wednesday)),
                this.getAppointmentsPerDay$(data.currentDate.clone().weekday(WeekDay.Thursday)),
                this.getAppointmentsPerDay$(data.currentDate.clone().weekday(WeekDay.Friday)),
                this.getAppointmentsPerDay$(data.currentDate.clone().weekday(WeekDay.Saturday)),
                this.getAppointmentsPerDay$(data.currentDate.clone().weekday(WeekDay.Sunday)),
                (
                    monday: IAppointmentModel[], 
                    tuesday: IAppointmentModel[], 
                    wednesday: IAppointmentModel[], 
                    thursday: IAppointmentModel[], 
                    friday: IAppointmentModel[],
                    saturday: IAppointmentModel[],
                    sunday: IAppointmentModel[]) => {

                        var appointmentsPerDay: { [day: string] : IAppointmentModel[] } = {};

                        appointmentsPerDay[WeekDay.Monday] = monday;
                        appointmentsPerDay[WeekDay.Tuesday] = tuesday;
                        appointmentsPerDay[WeekDay.Wednesday] = wednesday;
                        appointmentsPerDay[WeekDay.Thursday] = thursday;
                        appointmentsPerDay[WeekDay.Friday] = friday;
                        appointmentsPerDay[WeekDay.Saturday] = saturday;
                        appointmentsPerDay[WeekDay.Sunday] = sunday;

                        return appointmentsPerDay;
                    }
                );

            return zip.pipe(take(1));
        }
    }
}