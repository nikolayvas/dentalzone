import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { FormGroup, FormControl, AbstractControl } from '@angular/forms';
import { Subscription } from 'rxjs';

import { IScheduleRowModel, IDayOfWeekModel, IPaginatorModel, DayOrWeekMode, IAppointmentModel } from './schedule.models'

import * as moment from 'moment';
import { Utils } from '../../services/utils'
import { IPatientData } from '../../models/patient.dto';
import { DialogService } from 'primeng/api';
import { ChoosePatientComponent } from './choose-patient.component';
import { take } from 'rxjs/operators';

@Component({
    selector: 'schedule',
    templateUrl: './schedule.component.html',
    styleUrls: ['./schedule.component.css'],
    providers: [DialogService],
  })

export class ScheduleComponent implements OnInit, OnDestroy {
    private _subscription: Subscription = new Subscription();
    
    private _perMinutes: number = 15;
    private _startHour: number = 6;
    private _endHour: number = 19;
    private _daysInWeek: number = 7;

    private _currentMode: IPaginatorModel;
    private get columnsCount(): number {
        return this._currentMode.pageMode == DayOrWeekMode.Week ? this._daysInWeek : 1;
    }

    rows: IScheduleRowModel[] = [];
    cols: IDayOfWeekModel[] = [];

    form: FormGroup = new FormGroup({});

    constructor(
        public dialogService: DialogService
    ) { }

    ngOnInit() {
        
    }

    ngOnDestroy() {
        Utils.unsubscribe(this._subscription);
    }

    perionChanged(data: IPaginatorModel) {
        this._currentMode = data;

        this.initHeaders();
        this.initGridRows();
    }

    appointmentChanged(appointment: IAppointmentModel)
    {
        if(!appointment.hasPrev) {
            var upControlId = (appointment.y-1)*this._daysInWeek + appointment.x;
            var ctrl = this.form.controls[upControlId.toString()];

            ctrl.patchValue(Object.assign({}, ctrl.value, {hasPrev: true}));
        }
    }

    appointmentAdded(appointment: IAppointmentModel) {
        const ref = this.dialogService.open(ChoosePatientComponent, {
            header: 'Choose a patient',
            width: '70%',
            //contentStyle: {"max-height": "350px", "overflow": "auto"}
        });

        ref.onClose.pipe(take(1)).subscribe((patient: IPatientData) => {
            if (patient) {
                this.getContrloByPos(appointment.x, appointment.y).patchValue({"patientID": patient.id, "patientName": patient.firstName})
            }
        });
    }

    appointmentRemoved(appointment: IAppointmentModel) {
        this.form.controls
    }

    private getContrloByPos(x: number, y: number): AbstractControl {
        return this.form.controls[(y*this.columnsCount + x).toString()];
    }

    private initHeaders(): void {
        this.cols = [];

        if(this._currentMode.pageMode == DayOrWeekMode.Week) {
            for (let i = 1; i <= this._daysInWeek; i++) {
                const local = this._currentMode.currentDate.clone().weekday(i);
                this.cols.push( this.getColumnHeaderDate(local));
            }
        }
        else {
            const local = this._currentMode.currentDate.clone();
            this.cols.push( this.getColumnHeaderDate(local));
        }
    }
    
    private initGridRows() {
        Utils.unsubscribe(this._subscription);
        
        this.rows = [];
        this.removeFormControls();
        
        var newRows: IScheduleRowModel[] = [];
        var local = this._currentMode.currentDate.startOf('day');
        local.add('hours', this._startHour);

        var rowIndex = 0;
        while (local.hours() <= this._endHour) {

            const nextTime = local.format("HH:mm").toString();

            var appointmentsPerDay: IAppointmentModel[] = [];
            for (let i = 0; i <= this.columnsCount; i++) {
                const appointmentData = <IAppointmentModel>{x: i, y: rowIndex, }
                
                const control = new FormControl(appointmentData);
                this.form.addControl((rowIndex*this.columnsCount + i).toString(), control)
                
                this._subscription.add(control.valueChanges.subscribe(n=>{
                    this.appointmentChanged(n);
                }));
            }

            newRows.push({time: nextTime.endsWith("0") ? nextTime : "", forDays: appointmentsPerDay});

            local.add('minutes', this._perMinutes);
            rowIndex++;
        }

        this.rows = newRows;
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
}