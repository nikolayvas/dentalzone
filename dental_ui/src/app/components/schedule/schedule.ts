import { AbstractControl } from '@angular/forms';
import { IPaginatorModel, IAppointmentModel } from './schedule.models';
import { Utils } from '../../services/utils';
import { IAppointmentData } from '../../models/schedule-day.dto';

export class Schedule {
    static perMinutes: number = 15;
    static startHour: number = 6;
    static endHour: number = 19;
    static daysInWeek: number = 7;

    private _colorPerPatient: {[color: string]: string } = {};

    constructor(
        private controls: {[key: string] : AbstractControl}, 
        private cols: number, 
        private rows: number, 
        private mode: IPaginatorModel) {
    }

    addAppointment(appointment: IAppointmentModel) {
        var controlBelow = this.getControlByPos(appointment.x, appointment.y + 1);

        var control = this.getControlByPos(appointment.x, appointment.y);
        control.patchValue(Object.assign({}, appointment, {
            dateTime: this.getDateTimeByPos(appointment.x, appointment.y),
            patientID: appointment.patientID, 
            patientName: appointment.patientName,
            hasNext: !!controlBelow && !Utils.isBlankOrEmpty((<IAppointmentModel>controlBelow.value).patientID),
            color: this.getColorPerPatient(appointment.patientID)
        }));

        this.showExpandOfPreviousAppointment(appointment, false);
    }

    removeAppointment(appointment: IAppointmentModel) {
        var control = this.getControlByPos(appointment.x, appointment.y);
        control.patchValue(Object.assign({}, control.value, {
            patientID: undefined, 
            patientName: undefined, 
            hasNext: undefined,
            color: undefined
        }));

        this.showExpandOfPreviousAppointment(appointment, true);
    }

    extendAppointment(appointment: IAppointmentModel) {
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
            { hasNext: true }
        ));
    }

    getAppointmentsPerColumn(x: number): IAppointmentData[] {

        var res: IAppointmentData[] = [];

        for (let y = 0; y <= this.rows; y++) {
            const appointment = this.getControlByPos(x, y);
            if(appointment && appointment.value && appointment.value.patientID) {
                res.push(<IAppointmentData>{date: appointment.value.dateTime, patientID: appointment.value.patientID});
            }
        }

        return res;
    }

    getColumnDate(x: number): Date {
        var columnDate = this.mode.currentDate.clone().weekday(1);
        columnDate.add('day', x);

        return columnDate.toDate();
    }
    
    private getControlByPos(x: number, y: number): AbstractControl {
        return this.controls[this.getControlNameByPos(x, y)];
    }

    private getControlNameByPos(x: number, y: number): string {
        return (y*this.cols + x).toString();
    }

    private getDateTimeByPos(x: number, y: number): Date {
        var date = this.mode.currentDate.clone().weekday(1).startOf('day');
        date.add("days", x);
        date.add("hours", Schedule.startHour);
        date.add("minutes", Schedule.perMinutes*y);

        return date.toDate();
    }

    private showExpandOfPreviousAppointment(currentAppointment: IAppointmentModel, show: boolean): void {
        var controlUp = this.getControlByPos(currentAppointment.x, currentAppointment.y-1);
        if(controlUp) {
            controlUp.patchValue(Object.assign({}, controlUp.value, {"hasNext": !show}));
        }
    }

    private getColorPerPatient(patientId: string): string {
        if(!this._colorPerPatient[patientId]) {
            this._colorPerPatient[patientId] = this.getRandomColor();
        }

        return this._colorPerPatient[patientId];
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