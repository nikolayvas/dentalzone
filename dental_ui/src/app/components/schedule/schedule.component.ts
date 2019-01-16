import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { FormGroup, FormControl, AbstractControl } from '@angular/forms';
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
  })

export class ScheduleComponent {
    
    private _currentMode: IPaginatorModel;
    private _schedule: Schedule;

    private get columnsCount(): number {
        return this._currentMode.pageMode == DayOrWeekMode.Week ? Schedule.daysInWeek : 1;
    }

    rows: IScheduleRowModel[] = [];
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
            this.initGridRows();
            this.populateExistingsAppointments(appointments);
        });
    }

    appointmentAdded(appointment: IAppointmentModel) {
        const ref = this.dialogService.open(ChoosePatientComponent, {
            header: 'Choose a patient',
            width: '70%',
        });

        ref.onClose.pipe(take(1)).subscribe((patient: IPatientData) => {
            if (patient) {

                appointment.patientID = patient.id;
                appointment.patientName = patient.firstName;
               
                this._schedule.addAppointment(appointment);

                this.saveAppointmentsPerColumn(appointment.x);
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
                const dayOfWeek = this._currentMode.currentDate.clone().weekday(i);
                this.cols.push( this.getColumnHeaderDate(dayOfWeek));
            }
        }
        else {
            const singleDay = this._currentMode.currentDate.clone();
            this.cols.push( this.getColumnHeaderDate(singleDay));
        }
    }

    private initGridRows() {
        this.rows = [];
        this.removeFormControls();
        
        var newRows: IScheduleRowModel[] = [];
        var local = this._currentMode.currentDate.clone().startOf('day');
        local.add('hours', Schedule.startHour);

        var y = 0;
        while (local.hours() <= Schedule.endHour) {

            const nextTime = local.format("HH:mm").toString();

            for (let x = 0; x < this.columnsCount; x++) {
                const appointmentData = <IAppointmentModel>{x: x, y: y}
                
                const control = new FormControl(appointmentData);
                this.form.addControl(this.getControlNameByPos(x, y), control)
            }

            newRows.push({time: nextTime.endsWith("0") ? nextTime : ""});

            local.add('minutes', Schedule.perMinutes);
            y++;
        }

        this.rows = newRows;
        this._schedule = new Schedule(this.form.controls, this.columnsCount, newRows.length - 1, this._currentMode);
    }

    private populateExistingsAppointments(appointments: { [day: string] : IAppointmentModel[] }): void {
        for (let day in appointments) {
            const appointmentsList = appointments[day];
            for(let i = 0; i < appointmentsList.length; i++){
                const appointment = appointmentsList[i];

                var hoursAfterStartTime = appointment.dateTime.getHours() - Schedule.startHour;
                var minutes = appointment.dateTime.getMinutes();

                const y = hoursAfterStartTime*(60/Schedule.perMinutes) + minutes/Schedule.perMinutes;

                this._pService.getPatient$(appointment.patientID).pipe(take(1)).subscribe(p=>{
                    appointment.x = Number(day)-1;
                    appointment.y = y;
                    appointment.patientID= p.id;
                    appointment.patientName = p.firstName;

                    this._schedule.addAppointment(appointment);
                });
            }
        }
    }

    private getControlNameByPos(x: number, y: number): string {
        return (y*this.columnsCount + x).toString();
    }

    private removeFormControls() {
        for(var c in this.form.controls) {
            this.form.removeControl(c);
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