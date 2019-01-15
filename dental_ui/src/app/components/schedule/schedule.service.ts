import { Injectable } from '@angular/core';
import { Observable, zip as observableZip, of as observableOf } from 'rxjs';
import { take, filter, map } from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';

import { StoreService } from '../../services/store.service';
import { IAppointmentModel, IPaginatorModel, DayOrWeekMode, WeekDay } from './schedule.models';

import * as actions from '../../store/schedule.actions';

import * as moment from 'moment';
import { IAppointmentData, IScheduleByDayData } from '../../models/schedule-day.dto';
import { ScheduleByDayModel } from '../../models/schedule-day.model';

@Injectable()
export class ScheduleService {

    private _schedule$: Observable<ScheduleByDayModel[]>

    private get schedule$(): Observable<ScheduleByDayModel[]> {
        return this._schedule$ || (this._schedule$ = this._store.select(n=>n.data.clientPortalStore.scheduleState.perDay));
    }

    constructor(
        private http: HttpClient,
        private _store: StoreService
    ) {
    }

    getAppointmentsPerDay$(day: moment.Moment): Observable<IAppointmentModel[]> {

        this.appointmentsByDay$(day).pipe(take(1)).subscribe(n=> {
            if(!n) {
                this.seedAppointmentsPerDay(day);
            }
        });

        return this.appointmentsByDay$(day).pipe(filter(n=> !!n), map(n=>{
            return n.map(p=> {
               return <IAppointmentModel>{ patientID: p.patientID, dateTime: p.date } });
        }));
    }

    getAppointmentsPerDayOrWeek$(data: IPaginatorModel): Observable<{ [day: string] : IAppointmentModel[] }> {
        if(data.pageMode == DayOrWeekMode.Day) {
            return this.getAppointmentsPerDay$(data.currentDate.clone()).pipe(take(1), map(n=> {
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

    seedAppointmentsPerDay(day: moment.Moment): void {
        this.http.get<IAppointmentData[]>('/api/appointments', { params: { "day": day.toDate().toDateString() } }).pipe(
            take(1))
            .subscribe(data=> {
                const payload: IScheduleByDayData = <IScheduleByDayData>{day: day.toDate(), appointments: data || []};
                this._store.dispatch({ type: actions.ScheduleActions.SEED_DATA_FOR_DAY, payload: payload });
            });
    }

    saveAppointmentsPerDay(appointments: IAppointmentData[], day: Date): void {

        this.http.post('/api/appointments/update', {appointments: appointments, day: day}).pipe(
            take(1))
            .subscribe(
                _=> {
                    this._store.dispatch({ type: actions.ScheduleActions.UPDATE_DATA_FOR_DAY, payload: {appointments: appointments, day: day} });
                },
                err => console.log(err));
    }

    private appointmentsByDay$(day : moment.Moment): Observable<IAppointmentData[]> {
        return this.schedule$.pipe(
            map(n=>{
                const dayAsDate = day.toDate();
                return n.find(p=> 
                    p.day.getFullYear() === dayAsDate.getFullYear() &&
                    p.day.getMonth() === dayAsDate.getMonth() &&
                    p.day.getDate() === dayAsDate.getDate() 
                );
            }), 
            map(n=>{
                return n && n.appointments;
            }));
    }
}