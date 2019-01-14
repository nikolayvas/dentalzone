import { Component, forwardRef, Output, EventEmitter, ElementRef, Renderer2 } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';
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
    value: IAppointmentModel;

    @Output()
    onAdd = new EventEmitter<IAppointmentModel>();

    @Output()
    onRemove = new EventEmitter<IAppointmentModel>();

    @Output()
    onExtend = new EventEmitter<IAppointmentModel>();

    @Output()
    onShowInfo = new EventEmitter<IAppointmentModel>();

    constructor(
        private el: ElementRef,
        private renderer: Renderer2
        ) { }

    writeValue(appointment: IAppointmentModel): void {
        this.value = appointment;
        this.setBackGroundColorToParentTableCell();
    }

    registerOnChange(fn: any): void {
        this._onChangeCallback = fn;
    }

    registerOnTouched(fn: any): void {
    }

    extend() {
        this.onExtend.next(this.value);
    }

    add($event) {
        this.onAdd.next(this.value);
    }

    remove() {
        this.onRemove.next(this.value);
    }

    showInfo() {
        this.onShowInfo.next(this.value);
    }

    private emitValueChange() {
        this._onChangeCallback(this.value);
    }

    private setBackGroundColorToParentTableCell(): void {
        var color: string = this.value && this.value.color;

        if(!color) {
            if(this.value && this.value.x >= 5) {
                color = "beige";
            }
            else {
                color = "white";
            }
        }

        this.renderer.setStyle(this.el.nativeElement.parentNode, "background-color", color);
    }
}
