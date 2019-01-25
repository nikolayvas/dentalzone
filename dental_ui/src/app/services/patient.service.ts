
import {take, map, filter} from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { HttpClient, HttpParams } from '@angular/common/http';
import { IPatientData } from '../models/patient.dto';
import { StoreService } from './store.service'
import { GuidService } from './guid.service'
import * as actions from '../store/patient-actions';
import * as uiActions from '../store/ui-actions';
import { AuthService} from '../auth/auth.service'

@Injectable()
export class PatientService {

    private _patients$: Observable<IPatientData[]>;

    get patients$() : Observable<IPatientData[]> {
        return this._patients$ || (this._patients$ = this.store.select(n => n.data.clientPortalStore.patientsState.patients));
    }

    get patientSearchFilterChanged$(): Observable<string> {
        return this.store.select(n=>n.data.clientPortalStore.uiState.patients_filter);
    }

     constructor(
        private auth: AuthService,
        private http: HttpClient,
        private store: StoreService
    ) {
        auth.user$.subscribe(user => { 
            if(!!user) {
                this.seedPatientsList()
            }
        })
    }

    changeSearchFilter(filter: string): void {
        this.store.dispatch({ type: uiActions.UIActions.SET_PATIENTS_FILTER, payload: filter });
    }

    seedPatientsList(): void {
        this.http.get<IPatientData[]>('/api/patients').pipe(
            take(1))
            .subscribe(data=> {
                this.store.dispatch({ type: actions.PatientActions.SEED_PATIENTS, payload: data });
            });
    }

    getPatient$(id : string) : Observable<IPatientData> {
        return this.patients$.pipe(filter(n=> !!n), map(item => item.find(n=>n.id === id)))
    }

    updatePatientProfile(patient : IPatientData) : void {
        this.http.post('/api/patients/update', patient).pipe(
            take(1))
            .subscribe(
                data=> {
                    this.store.dispatch({ type: actions.PatientActions.UPDATE_PATIENT, payload: patient });
                },
                err => console.log(err));
    }

    addNewPatientProfile(patient : IPatientData) : void {

        const payload = Object.assign({}, patient, { id: GuidService.newGuid() })

        this.http.post('/api/patients/create', payload).pipe(
        take(1))
        .subscribe(
            patientId=> {
                const newPayload = Object.assign({}, payload, {id: patientId});
                this.store.dispatch({ type: actions.PatientActions.CREATE_PATIENT, payload: newPayload });
            },
            err => console.log(err));
    }

    removePatient(id : string) : void {
        this.http.delete('/api/patients/remove', {params: {"id": id}}).pipe(
        take(1))
        .subscribe(
            data=> {
                this.store.dispatch({ type: actions.PatientActions.REMOVE_PATIENT, payload: id });
            },
            err => console.log(err));
    }
}