
import { Subscription } from 'rxjs';
import { Component, OnInit, OnDestroy, ChangeDetectionStrategy, ChangeDetectorRef } from '@angular/core';

import { MenuItem } from 'primeng/primeng';
import { PatientService } from '../../services/patient.service';
import { Utils } from '../../services/utils';
import { filter } from 'rxjs/operators';

@Component({
    selector: 'app-menu',
    templateUrl: 'menu.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class MenuComponent implements OnInit, OnDestroy  {
    private subscriptions: Subscription = new Subscription();

    items: MenuItem[];

    constructor(
        private changeDedectionRef: ChangeDetectorRef,
        private patientService: PatientService
    ) {
    }

    ngOnInit() {
        this.subscriptions.add(this.patientService.currentParient$.pipe(filter(n=>!!n)).subscribe(patient=>{
            
            const link1 = "patients/edit-patient-profile/" + patient.id;
            const link2 = "patients/tooth-status/" + patient.id;
            const link3 = "patients/upload/" + patient.id;
            const link4 = "patients/download/" + patient.id;

            this.items = [
                { label: 'Patients', icon: 'icon-user-alt', routerLink: 'patients', command: (event) => {} },
                { label: 'Schedule', icon: 'icon-calendar', routerLink: 'schedule', command: (event) => {} },
                { label: patient.name, icon: 'icon-user', command: (event) => {} , items: [
                    { label: 'Profile', routerLink: link1, command: (event) => {}  },
                    { label: 'Tooth Status', routerLink: link2, command: (event) => {} },
                    { label: 'File Upload', routerLink: link3, command: (event) => {} },
                    { label: 'File Download', routerLink: link4, command: (event) => {} },
                ] },
            ];

            this.changeDedectionRef.detectChanges();
        }));

        this.items = [
            { label: 'Patients', icon: 'icon-user-alt', routerLink: 'patients', command: (event) => {} /*this.disableAllBut(event.item)*/ },
            { label: 'Schedule', icon: 'icon-calendar', routerLink: 'schedule', command: (event) => {} /*this.disableAllBut(event.item)*/ },
        ];
    }

    ngOnDestroy() {
        Utils.unsubscribe(this.subscriptions);
    }

    disableAllBut(except: any) {
        for (let item of this.items) {
            if (item.label !== except.label) {
                item.expanded = false;
            }
        }
    }
}