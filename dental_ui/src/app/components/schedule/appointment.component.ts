import { Component, Input, Output, EventEmitter, ElementRef, Renderer2, ChangeDetectionStrategy, OnInit, OnDestroy } from '@angular/core';
import { IAppointmentModel } from './schedule.models';

@Component({
    selector: 'appointment',
    templateUrl: './appointment.component.html',
    styleUrls: ['appointment.component.scss'],
    //changeDetection: ChangeDetectionStrategy.OnPush
})
export class AppointmentComponent implements OnInit, OnDestroy {
    @Input() 
    appointment: IAppointmentModel;

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

    ngOnInit() {
        this.setBackGroundColorToParentTableCell();
    }

    ngOnDestroy() {
        this.undoBackGroundColorToParentTableCell();
    }

    extend() {
        this.onExtend.next(this.appointment);
    }

    remove() {
        this.onRemove.next(this.appointment);
    }

    showInfo() {
        this.onShowInfo.next(this.appointment);
    }

    private setBackGroundColorToParentTableCell(): void {
        var color: string = this.appointment.color;

        if(!color) {
            if(this.appointment.x >= 5) {
                color = "beige";
            }
            else {
                color = "white";
            }
        }

        this.renderer.setStyle(this.el.nativeElement.parentNode, "background-color", color);
    }

    undoBackGroundColorToParentTableCell(): void {
        this.renderer.setStyle(this.el.nativeElement.parentNode, "background-color", "white");
    }
}
