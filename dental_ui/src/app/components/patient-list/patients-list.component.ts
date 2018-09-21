import { Component, OnInit, OnDestroy } from "@angular/core";
import { Observable, Subscription } from 'rxjs';

import { IPatientData } from '../../models/patient.dto';
import { PatientService } from '../../services/patient.service'
import { ConfirmationService } from 'primeng/primeng';

import { Utils } from '../../services/utils'

@Component({
    selector: 'patients-list',
    templateUrl: './patients-list.component.html'
})

export class PatientsListComponent implements OnInit, OnDestroy {

    private subscriptions: Subscription = new Subscription();
    patients$: Observable<IPatientData[]>;
    searchFor: string;

    constructor(
        private service: PatientService ,
        private confirmationService: ConfirmationService,
    ) {
        this.patients$ = this.service.patients$;
        this.subscriptions.add(this.service.patientSearchFilterChanged$.subscribe(n=>this.searchFor = n));
    }

    ngOnInit() {
       
    }

    ngOnDestroy() {
        Utils.unsubscribe(this.subscriptions);
    }

    inputChange(control) {
        const filter = control.target.value;
        this.service.changeSearchFilter(filter);
    }

    remove(patient : IPatientData) {
        this.confirmationService.confirm({
            header: 'Remove patient',
            message: 'Are you sure you want to remove that patient?',
            icon: 'fa fa-question-circle',
            accept: () => {
                this.service.removePatient(patient.id);
            }
        });
    }
}