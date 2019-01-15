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
import { Utils } from '../../services/utils'

import * as moment from 'moment';
import { IAppointmentData } from '../../models/schedule-day.dto';

@Component({
    selector: 'schedule',
    templateUrl: './schedule.component.html',
    styleUrls: ['./schedule.component.css'],
    providers: [DialogService],
  })

export class ScheduleComponent {
    
    private _perMinutes: number = 15;
    private _startHour: number = 6;
    private _endHour: number = 19;
    private _daysInWeek: number = 7;

    private _currentMode: IPaginatorModel;
    private _colorPerPatient: {[color: string]: string } = {};

    private get columnsCount(): number {
        return this._currentMode.pageMode == DayOrWeekMode.Week ? this._daysInWeek : 1;
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

        this._service.getAppointmentsPerDayOrWeek$(data).subscribe(n=> {
            this.populateGrid(n);
        });
    }

    appointmentExtended(appointment: IAppointmentModel)
    {
        var controlBelow = this.getControlByPos(appointment.x, appointment.y + 1);
        var nextControl = this.getControlByPos(appointment.x, appointment.y + 2);

        controlBelow.patchValue(Object.assign({}, controlBelow.value, {
            dateTime: this.getDateTimeByPos(appointment.x, appointment.y + 1),
            patientID: appointment.patientID,
            patientName: appointment.patientName,
            color: appointment.color,
            hasNext: nextControl && nextControl.value.patientID,
        }));

        var control = this.getControlByPos(appointment.x, appointment.y);
        control.patchValue(Object.assign({}, control.value, 
            { hasNext: true }));

        this._service.saveAppointmentsPerDay(this.getAppointmentsPerColumn(appointment.x), this.getDatePerColumn(appointment.x));
    }

    appointmentAdded(appointment: IAppointmentModel) {
        const ref = this.dialogService.open(ChoosePatientComponent, {
            header: 'Choose a patient',
            width: '70%',
            //contentStyle: {"max-height": "350px", "overflow": "auto"}
        });

        ref.onClose.pipe(take(1)).subscribe((patient: IPatientData) => {
            if (patient) {
                var controlBelow = this.getControlByPos(appointment.x, appointment.y + 1);

                var control = this.getControlByPos(appointment.x, appointment.y);
                control.patchValue(Object.assign({}, appointment, {
                    dateTime: this.getDateTimeByPos(appointment.x, appointment.y),
                    patientID: patient.id, 
                    patientName: patient.firstName,
                    hasNext: !!controlBelow && !Utils.isBlankOrEmpty((<IAppointmentModel>controlBelow.value).patientID),
                    color: this.getColorPerPatient(patient.id)
                }));

                this.showExpandOfPreviousAppointment(appointment, false);

                this._service.saveAppointmentsPerDay(this.getAppointmentsPerColumn(appointment.x), this.getDatePerColumn(appointment.x));
            }
        });
    }

    appointmentRemoved(appointment: IAppointmentModel) {
        var control = this.getControlByPos(appointment.x, appointment.y);
        control.patchValue(Object.assign({}, control.value, {
            patientID: undefined, 
            patientName: undefined, 
            hasNext: undefined,
            color: undefined
        }));

        this.showExpandOfPreviousAppointment(appointment, true);

        this._service.saveAppointmentsPerDay(this.getAppointmentsPerColumn(appointment.x), this.getDatePerColumn(appointment.x));
    }

    showPatientInfo(appointment: IAppointmentModel) {
        this._pService.changeSearchFilter(appointment.patientID);
        this.router.navigateByUrl('app/portal/patients');
    }

    private populateGrid(appointments: { [day: string] : IAppointmentModel[] }): void {

        this.initHeaders();
        this.initGridRows();

        for (let day in appointments) {

            const appointmentsList = appointments[day];
            for(let i = 0; i<appointmentsList.length; i++){
                const appointment = appointmentsList[i];

                var hoursAfterStartTime = appointment.dateTime.getHours() - this._startHour;
                var minutes = appointment.dateTime.getMinutes();

                const y = hoursAfterStartTime*(60/this._perMinutes) + minutes/this._perMinutes;

                this._pService.getPatient$(appointment.patientID).pipe(take(1)).subscribe(p=>{
                    appointment.x = Number(day)-1;
                    appointment.y = y;
                    appointment.color = this.getColorPerPatient(appointment.patientID);
                    appointment.patientName = p.firstName;

                    this.getControlByPos(appointment.x, appointment.y).patchValue(appointment);
                    this.showExpandOfPreviousAppointment(appointment, false);
                });
            }
        }
    }

    private getAppointmentsPerColumn(x: number): IAppointmentData[] {

        var res: IAppointmentData[] = [];

        for (let y = 0; y < this.rows.length; y++) {
            const appointment = this.getControlByPos(x, y);
            if(appointment && appointment.value && appointment.value.patientID) {
                res.push(<IAppointmentData>{date: appointment.value.dateTime, patientID: appointment.value.patientID});
            }
        }

        return res;
    }

    private getDatePerColumn(x: number): Date {
        var columnDate = this._currentMode.currentDate.clone().weekday(1);
        columnDate.add('day', x);

        return columnDate.toDate();
    }

    private initHeaders(): void {
        this.cols = [];

        if(this._currentMode.pageMode == DayOrWeekMode.Week) {
            for (let i = 1; i <= this._daysInWeek; i++) {
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
        local.add('hours', this._startHour);

        var y = 0;
        while (local.hours() <= this._endHour) {

            const nextTime = local.format("HH:mm").toString();

            for (let x = 0; x < this.columnsCount; x++) {
                const appointmentData = <IAppointmentModel>{x: x, y: y}
                
                const control = new FormControl(appointmentData);
                this.form.addControl(this.getControlNameByPos(x, y), control)
            }

            newRows.push({time: nextTime.endsWith("0") ? nextTime : ""});

            local.add('minutes', this._perMinutes);
            y++;
        }

        this.rows = newRows;
    }

    private getColorPerPatient(patientId: string): string {
        if(!this._colorPerPatient[patientId]) {
            this._colorPerPatient[patientId] = this.getRandomColor();
        }

        return this._colorPerPatient[patientId];
    }

    private getDateTimeByPos(x: number, y: number): Date {
        var date = this._currentMode.currentDate.clone().weekday(1).startOf('day');
        date.add("days", x);
        date.add("hours", this._startHour);
        date.add("minutes", this._perMinutes*y);

        return date.toDate();
    }

    private getControlByPos(x: number, y: number): AbstractControl {
        return this.form.controls[this.getControlNameByPos(x, y)];
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
            isToday: today.dayOfYear() === date.dayOfYear()
        }
    }

    private showExpandOfPreviousAppointment(currentAppointment: IAppointmentModel, show: boolean): void {
        var controlUp = this.getControlByPos(currentAppointment.x, currentAppointment.y-1);
        if(controlUp) {
            controlUp.patchValue(Object.assign({}, controlUp.value, {"hasNext": !show}));
        }
    }

    private randColor() {
        var color = (function lol(m, s, c) {
                        return s[m.floor(m.random() * s.length)] +
                            (c && lol(m, s, c - 1));
                    })(Math, '3456789ABCDEF', 4);
        return color;
    }

    private getRandomColor(): string {
        const color = this.randColor();
        return '#' + color;

        /*
        var color = Math.floor(0x1000000 * Math.random()).toString(16);
        return '#' + ('000000' + color).slice(-6);
        */
    }
}