import { Component, OnDestroy } from "@angular/core";
import { FormGroup, FormControl } from '@angular/forms';
import { Observable, Subscription } from 'rxjs';

import { IPatientData } from '../../models/patient.dto';
import { PatientService } from '../../services/patient.service'
import { ConfirmationService } from 'primeng/primeng';

import { Utils } from '../../services/utils'

@Component({
    selector: 'patients-list',
    templateUrl: './patients-list.component.html'
})

export class PatientsListComponent implements OnDestroy {

    form: FormGroup;

    search: FormControl;

    private subscriptions: Subscription = new Subscription();
    patients$: Observable<IPatientData[]>;
    searchFor: string;

    constructor(
        private service: PatientService ,
        private confirmationService: ConfirmationService,
    ) {
        this.form = new FormGroup({
            search: this.search = new FormControl('')
        });

        this.patients$ = this.service.patients$;

        this.initDataSubscriptions();
    }

    ngOnDestroy() {
        Utils.unsubscribe(this.subscriptions);
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

    private initDataSubscriptions() {
        this.subscriptions.add(this.search.valueChanges.subscribe(n=> {
            this.service.changeSearchFilter(n);
        }));

        this.subscriptions.add(this.service.patientSearchFilterChanged$.subscribe(n=>{
            this.search.patchValue(n, { emitEvent: false });
            this.searchFor = n;
        }));
    }
}