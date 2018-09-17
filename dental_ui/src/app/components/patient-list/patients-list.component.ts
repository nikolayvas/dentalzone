import { Component, OnInit, OnDestroy } from "@angular/core";
import { Observable } from 'rxjs';
import { Router } from '@angular/router';

import { IPatientData } from '../../models/patient.dto';
import { PatientService } from '../../services/patient.service'
import { Utils } from '../../services/utils'

@Component({
    selector: 'patients-list',
    templateUrl: './patients-list.component.html'
})

export class PatientsListComponent implements OnInit, OnDestroy {

    patients$: Observable<IPatientData[]>;
    private searchFor: string;

    constructor(
        private router: Router,
        private service: PatientService 
    ) {
        this.patients$ = this.service.patients$;
    }

    ngOnInit() {
       
    }

    ngOnDestroy() {
        
    }

    inputChange(control) {
        this.searchFor = control.target.value;
    }

    remove(patient : IPatientData) {
        this.service.removePatient(patient.id);
    }
}