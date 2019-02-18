import { Component, ChangeDetectionStrategy } from '@angular/core';
import { Router } from '@angular/router';
import { FormGroup } from '@angular/forms';
import { take,  } from 'rxjs/operators';

import { IScheduleRowModel, IDayOfWeekModel, IPaginatorModel, DayOrWeekMode, IAppointmentModel } from './schedule.models'
import { IPatientData } from '../../models/patient.dto';
import { DialogService } from 'primeng/api';
import { ChoosePatientComponent } from './choose-patient.component';
import { ScheduleService } from './schedule.service';
import { PatientService } from '../../services/patient.service';

import * as moment from 'moment';
import { Schedule } from './schedule';

@Component({
    selector: 'schedule',
    templateUrl: './schedule.component.html',
    styleUrls: ['./schedule.component.css'],
    providers: [DialogService],
    //changeDetection: ChangeDetectionStrategy.OnPush
  })

export class ScheduleComponent {
    counter: number = 0;

    private _currentMode: IPaginatorModel;
    private _schedule: Schedule;

    /*
    get runChangeDetection() {
        console.log('checking the view' + this.counter++);
        return true;
    }
    */

    get rows(): IScheduleRowModel[] {
        return this._schedule ? this._schedule.rows : [];
    }
    
    cols: IDayOfWeekModel[] = [];

    form: FormGroup = new FormGroup({});

    constructor(
        public dialogService: DialogService,
        private _service: ScheduleService,
        private _pService: PatientService,
        private router: Router,
    ) { }

    periodChanged(data: IPaginatorModel) {
        this._currentMode = data;

        this._service.getAppointmentsPerDayOrWeek$(data).subscribe(appointments => {
            this.initHeaders();
            this.initAppointmentsData(appointments);
        });
    }

    getAppointment(x: number, y: number): IAppointmentModel {
        return this._schedule.getAppointment(x, y);
    }

   addAppointment(x: number, y: number) {
        const ref = this.dialogService.open(ChoosePatientComponent, {
            header: 'Choose a patient',
            width: '70%',
        });

        ref.onClose.pipe(take(1)).subscribe((patient: IPatientData) => {
            if (patient) {
                this._schedule.addAppointment(x, y, patient);
                this.saveAppointmentsPerColumn(x);
            }
        });
    }

    appointmentExtended(appointment: IAppointmentModel)
    {
        this._schedule.extendAppointment(appointment);
        this.saveAppointmentsPerColumn(appointment.x);
    }

    appointmentRemoved(appointment: IAppointmentModel) {
        this._schedule.removeAppointment(appointment);
        this.saveAppointmentsPerColumn(appointment.x);
    }

    showPatientInfo(appointment: IAppointmentModel) {
        this._pService.changeSearchFilter(appointment.patientID);
        this.router.navigateByUrl('app/portal/patients');
    }

    private saveAppointmentsPerColumn(x: number) {
        this._service.saveAppointmentsPerDay(
            this._schedule.getAppointmentsPerColumn(x), 
            this._schedule.getColumnDate(x)
        );
    }

    private initHeaders(): void {
        this.cols = [];

        if(this._currentMode.pageMode == DayOrWeekMode.Week) {
            for (let i = 1; i <= Schedule.daysInWeek; i++) {
                const dayOfWeek = this._currentMode.currentDate.clone().isoWeekday(i);
                this.cols.push( this.getColumnHeaderDate(dayOfWeek));
            }
        }
        else {
            const singleDay = this._currentMode.currentDate.clone();
            this.cols.push( this.getColumnHeaderDate(singleDay));
        }
    }

    private initAppointmentsData(appointments: { [day: string] : IAppointmentModel[] }): void {

        this._schedule = new Schedule(this._currentMode);

        for (let day in appointments) {
            const appointmentsList = appointments[day];
            for(let i = 0; i < appointmentsList.length; i++){
                const appointment = appointmentsList[i];

                var hoursAfterStartTime = appointment.dateTime.getHours() - Schedule.startHour;
                var minutes = appointment.dateTime.getMinutes();

                const y = hoursAfterStartTime*(60/Schedule.perMinutes) + minutes/Schedule.perMinutes;

                this._pService.getPatient$(appointment.patientID).pipe(take(1)).subscribe(p=>{
                    if(p) {
                        this._schedule.addAppointment(Number(day)-1, y, p);
                    }
                });
            }
        }
    }

    private getColumnHeaderDate(date: moment.Moment): IDayOfWeekModel {
        const today = moment();

        return <IDayOfWeekModel>{
            dayOfWeek: date.format('dddd'),
            dayOfMonth: date.format('MMM DD'),
            isToday: today.dayOfYear() === date.dayOfYear(),
        }
    }
}