import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';
import { ActivatedRoute, Router } from "@angular/router";
import { Subscription } from 'rxjs';
import { Utils } from '../../services/utils'

import { InputValidators } from '../../validation/input-validators'
import { PatientService } from '../../services/patient.service'
import { IPatientData } from '../../models/patient.dto';
import { DatePipe } from '@angular/common';


@Component({
    selector: 'add-edit-patient-profile',
    templateUrl: './add-edit-patient.component.html'
})
export class PatientProfileComponent implements OnInit, OnDestroy {

    patientProfileForm: FormGroup;

    firstName: FormControl;
    middleName: FormControl;
    lastName: FormControl;
    email: FormControl;
    address: FormControl;
    phoneNumber: FormControl;
    generalInfo: FormControl;
    registrationDate: FormControl;

    private patientId: string;
    private isNew: boolean = false;

    private subscriptions: Subscription = new Subscription();

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        private facade: PatientService,
        private datePipe: DatePipe,
        private patientService: PatientService
    ) {

    }

    ngOnInit() {

        this.patientId = this.route.snapshot.params['id'];

        if (!this.patientId) {
            this.isNew = true;
        }

        this.patientProfileForm = new FormGroup({
            firstName: this.firstName = new FormControl(''),
            middleName: this.middleName = new FormControl(''),
            lastName: this.lastName = new FormControl(''),
            email: this.email = new FormControl('', InputValidators.validateEmail),
            address: this.address = new FormControl(''),
            phoneNumber: this.phoneNumber = new FormControl(''),
            generalInfo: this.generalInfo = new FormControl(''),
            registrationDate: this.registrationDate = new FormControl({value: '', disabled: true}),
        });

        this.initDataSubscriptions();
    }

    private initDataSubscriptions() {
        if (!this.isNew) {
            this.subscriptions.add(this.facade.getPatient$(this.patientId).subscribe(patient => {
                if (!!patient) {
                    this.patientService.currentPatientHasChanged(patient.id, patient.firstName);
                    const formatted = Object.assign({}, patient, {registrationDate: this.datePipe.transform(patient.registrationDate,"dd-MM-yyyy")});
                    this.patientProfileForm.patchValue(formatted);
                }
            }));
        }
    }

    ngOnDestroy() {
        Utils.unsubscribe(this.subscriptions);
    }

    submit() {

        let patientData: IPatientData = Object.assign({}, this.patientProfileForm.value, { id: this.patientId });

        if(!this.isNew) {
            this.facade.updatePatientProfile(patientData);
        }
        else {
            this.facade.addNewPatientProfile(patientData);
        }

        this.router.navigateByUrl('/app/portal/patients');
    }

    back() {
        this.router.navigateByUrl('/app/portal/patients');  
    }
}