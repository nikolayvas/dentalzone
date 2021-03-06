import {take, map, filter} from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { Observable, BehaviorSubject, of } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { IPatientData } from '../models/patient.dto';
import { StoreService } from './store.service'
import { GuidService } from './guid.service'
import * as actions from '../store/patient-actions';
import * as uiActions from '../store/ui-actions';
import { AuthService} from '../auth/auth.service'

@Injectable()
export class PatientService {

    private _patients$: Observable<IPatientData[]>;

    private _currentParient$: BehaviorSubject<{id: string, name: string}> = new BehaviorSubject<{id: string, name: string}>(undefined);

    get patients$() : Observable<IPatientData[]> {
        return this._patients$ || (this._patients$ = this.store.select(n => n.data.clientPortalStore.patientsState.patients));
    }

    get patientSearchFilterChanged$(): Observable<string> {
        return this.store.select(n=>n.data.clientPortalStore.uiState.patients_filter);
    }

    get currentParient$(): Observable<{id: string, name: string}> {
        return this._currentParient$;
    }

     constructor(
        private auth: AuthService,
        private http: HttpClient,
        private store: StoreService,
    ) {
        auth.user$.subscribe(user => { 
            if(!!user) {
                this.seedPatientsList()
            }
        })
    }

    currentPatientHasChanged(id: string, name: string): void {
        if(!this._currentParient$.value || (this._currentParient$.value && this._currentParient$.value.id != id)) {
            this._currentParient$.next({id: id, name: name})
        }
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

    getTagsPerPatient(id: string): Observable<string[]> {

        //return of(["one", "two", "two", "two", "two", "two", "two", "two", "two", "two", "two"]);
        return this.http.get('/api/patient/filesTags', {params: { "id": id }}).pipe(take(1), map(n =>{ return n as string[] }));
    }
}