import { IPaginatorModel, IAppointmentModel, IScheduleRowModel, DayOrWeekMode } from './schedule.models';
import { IAppointmentData } from '../../models/schedule-day.dto';
import { Color } from './color';
import { IPatientData } from '../../models/patient.dto';

export class Schedule {

    static perMinutes: number = 15;
    static startHour: number = 6;
    static endHour: number = 14;
    static daysInWeek: number = 7;

    private _colorPerPatient: {[color: string]: string } = {};
    private _appointments: IAppointmentModel[];
    private _rows: IScheduleRowModel[];

    get rows(): IScheduleRowModel[]{
        return this._rows || (this._rows = this.getRows());
    }

    constructor(private mode: IPaginatorModel) {
        this._appointments = new Array(mode.pageMode*this.rows.length).fill(null);
    }

    addAppointment(x: number, y: number, patient: IPatientData) {
        let appointment = <IAppointmentModel> {
            x: x, 
            y: y,
            dateTime: this.getDateTimeByPos(x, y),
            patientID: patient.id,
            patientName: patient.firstName,
            color: this.getColorPerPatient(patient.id),
            hasNext: this.hasNext(x, y)
        }

        this._appointments[y*this.mode.pageMode + x] = appointment;
        this.showExpandOfPreviousAppointment(appointment, false);
    }

    getAppointment(x: number, y: number): IAppointmentModel {
        return this._appointments[y*this.mode.pageMode + x];
    }
  
    removeAppointment(appointment: IAppointmentModel) {
        this._appointments[appointment.y*this.mode.pageMode + appointment.x] = null;
        
        this.showExpandOfPreviousAppointment(appointment, true);
    }

    extendAppointment(appointment: IAppointmentModel) {
        let actorAppointment = this.getAppointment(appointment.x, appointment.y);
        actorAppointment.hasNext = true;

        let appointment3 = this.getAppointment(appointment.x, appointment.y + 2);
        let newAppointment = Object.assign({}, actorAppointment, {
            x: appointment.x,
            y: appointment.y+1,
            dateTime: this.getDateTimeByPos(appointment.x, appointment.y + 1), 
            hasNext: appointment3 !== null 
        });

        this._appointments[(appointment.y+1)*this.mode.pageMode + appointment.x] = newAppointment;
    }

    getAppointmentsPerColumn(x: number): IAppointmentData[] {

        var res: IAppointmentData[] = [];

        for (let y = 0; y <= this.rows.length; y++) {
            const appointment = this.getAppointment(x, y);
            if(appointment) {
                res.push(<IAppointmentData>{
                    date: appointment.dateTime, 
                    patientID: appointment.patientID
                });
            }
        }

        return res;
    }

    getColumnDate(x: number): Date {
        var columnDate = this.mode.currentDate.clone().isoWeekday(1);
        columnDate.add('day', x);

        return columnDate.toDate();
    }
    
    private getDateTimeByPos(x: number, y: number): Date {
        var date = this.mode.currentDate.clone().isoWeekday(1).startOf('day');
        date.add("days", x);
        date.add("hours", Schedule.startHour);
        date.add("minutes", Schedule.perMinutes*y);

        return date.toDate();
    }

    private showExpandOfPreviousAppointment(currentAppointment: IAppointmentModel, show: boolean): void {
        var controlUp = this.getAppointment(currentAppointment.x, currentAppointment.y-1);
        if(controlUp) {
            controlUp.hasNext= !show;
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

        return "#" + Color.Colors[this.randomIntFromInterval(0, Color.Colors.length -1)];
        /*
        const color = this.randColor();
        return '#' + color;
        */

        /*
        var color = Math.floor(0x1000000 * Math.random()).toString(16);
        return '#' + ('000000' + color).slice(-6);
        */
    }

    private randomIntFromInterval(min: number, max: number): number
    {
        return Math.floor(Math.random()*(max-min+1)+min);
    }

    private getRows(): IScheduleRowModel[] {

        var newRows: IScheduleRowModel[] = [];
        var local = this.mode.currentDate.clone().startOf('day');
        local.add('hours', Schedule.startHour);

        var y = 0;
        while (local.hours() <= Schedule.endHour) {

            const nextTime = local.format("HH:mm").toString();

            newRows.push({time: nextTime.endsWith("0") ? nextTime : ""});

            local.add('minutes', Schedule.perMinutes);
            y++;
        }

        return newRows;
    }

    private hasNext(x: number, y: number): boolean {
        let appointmentIndex = (y+1)*this.mode.pageMode + x;
        if(this._appointments.length < appointmentIndex) {
            return false;
        }
        else {
            return this._appointments[appointmentIndex] !== null;
        }
    }
}