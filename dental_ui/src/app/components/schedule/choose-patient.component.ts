import { Component, OnDestroy, ViewChild, Output, EventEmitter } from "@angular/core";
import { FormGroup, FormControl } from '@angular/forms';
import { Observable, Subscription } from 'rxjs';

import { IPatientData } from '../../models/patient.dto';
import { PatientService } from '../../services/patient.service'
import { DynamicDialogRef, DynamicDialogConfig } from 'primeng/primeng';

import { Utils } from '../../services/utils'
import { Table } from 'primeng/table';

@Component({
    selector: 'choose-patient',
    templateUrl: './choose-patient.component.html'
})

export class ChoosePatientComponent implements OnDestroy {

    form: FormGroup;

    search: FormControl;

    private subscriptions: Subscription = new Subscription();
    patients$: Observable<IPatientData[]>;
    searchFor: string;

    @ViewChild('dt') table: Table;

    constructor(
        private service: PatientService,
        public ref: DynamicDialogRef, public config: DynamicDialogConfig
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

    select(patient : IPatientData) {
        this.ref.close(patient);
    }

    cancel() {
        this.ref.close(undefined);
    }

    private initDataSubscriptions() {
        this.subscriptions.add(this.search.valueChanges.subscribe(n=> {
            this.service.changeSearchFilter(n);
        }));

        this.subscriptions.add(this.service.patientSearchFilterChanged$.subscribe(n=>{
            this.search.patchValue(n, { emitEvent: false });
            this.searchFor = n;

            if(this.table) {
                this.table.first=0;
            }
        }));
    }
}