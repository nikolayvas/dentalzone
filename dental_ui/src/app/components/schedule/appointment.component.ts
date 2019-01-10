import { Component, forwardRef, Output, EventEmitter } from '@angular/core';
import {
    FormGroup, FormControl,
    ControlValueAccessor, NG_VALUE_ACCESSOR
} from '@angular/forms';
import { Utils } from '../../services/utils'
import { IAppointmentModel } from './schedule.models';

const APPOINTMENT_CONTROL_VALUE_ACCESSOR = {
    provide: NG_VALUE_ACCESSOR,
    useExisting: forwardRef(() => AppointmentComponent),
    multi: true
};

@Component({
    selector: 'appointment',
    templateUrl: './appointment.component.html',
    providers: [APPOINTMENT_CONTROL_VALUE_ACCESSOR],
    styleUrls: ['appointment.component.scss']
})
export class AppointmentComponent implements ControlValueAccessor {
    private _onChangeCallback: (_: any) => void = Utils.noop;

    @Output()
    onAdd = new EventEmitter<IAppointmentModel>();

    @Output()
    onRemove = new EventEmitter<IAppointmentModel>();

    form: FormGroup;

    hasPrev: FormControl;
    hasNext: FormControl;
    patientID: FormControl;
    patientName: FormControl;

    constructor(
        ) {
            this.form = new FormGroup({
                x: new FormControl(''),
                y: new FormControl(''),
                hasPrev: this.hasPrev = new FormControl(false),
                hasNext: this.hasNext = new FormControl(false),
                patientID: this.patientID = new FormControl(''),
                patientName: this.patientName = new FormControl(''),
            });
        }

    writeValue(appointment: IAppointmentModel): void {
        this.form.patchValue(appointment, {emitEvent: false});
    }

    registerOnChange(fn: any): void {
        this._onChangeCallback = fn;
    }

    registerOnTouched(fn: any): void {
    }

    clickUp() {
        this.hasPrev.patchValue(true);
        this.emitValueChange();
    }

    clickDown() {
        this.hasNext.patchValue(true);
        this.emitValueChange();
    }

    addItem($event) {
        this.onAdd.next(this.form.value)
    }

    removeItem() {
        this.onRemove.next(this.form.value)
    }

    private emitValueChange() {
        this._onChangeCallback(this.form.value);
    }
}
