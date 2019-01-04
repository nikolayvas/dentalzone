
import {take, map} from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { HttpClient, HttpParams } from '@angular/common/http';
import { StoreService } from './store.service'
import { GuidService } from './guid.service'
import * as actions from '../store/teeth-data.action';
import { IToothActionData } from '../models/tooth-action.dto'
import { ITeethData } from '../models/teeth-data.dto'
import { ToothActionModel } from '../models/tooth-action.model'

@Injectable()
export class ToothStatusService {

    constructor(
        private http: HttpClient,
        private store: StoreService
    ) {
        
    }

    toothDiagnosisData$(toothNo?: string): Observable<ToothActionModel[]> {

        const result = this.store.select(n => n.data.clientPortalStore.teethDataState.toothDiagnosis);

        if(!toothNo) {
            return result;
        }
        else {
            return result.pipe(map(n=>n.filter(x=>x.toothNo == toothNo)));
        }
    }

    toothManipulationsData$(toothNo?: string): Observable<ToothActionModel[]> {
        const result = this.store.select(n => n.data.clientPortalStore.teethDataState.toothManipupations);

        if(!toothNo) {
            return result;
        }
        else {
            return result.pipe(map(n=>n.filter(x=>x.toothNo == toothNo)));
        }
    }

    seedTeethData(patientId: string) : void {
     
        this.http.get<ITeethData>('/api/toothStatus/seedTeethData', { params: { "patientId": patientId } }).pipe(
            take(1))
            .subscribe(
                data=> {
                    this.store.dispatch({ type: actions.TeethDataActions.SEED_TEETH_STATE_DATA, payload: data });
                },
                err => console.log(err));
    }

    clearTeethData(): void {
        this.store.dispatch({ type: actions.TeethDataActions.CLEAR_TEETH_STATE_DATA });
    }

    addDiagnosis(toothNo: string, diagnosisId: number, patientId: string) : void {

        var diagnosis = <IToothActionData>{
            id: GuidService.newGuid(),
            patientId:patientId, 
            toothNo: toothNo, 
            actionId: Number(diagnosisId),
            date: new Date()
        }
        
        this.http.post('/api/toothStatus/addToothDiagnosis', diagnosis).pipe(
            take(1))
            .subscribe(
                _=> {
                    this.store.dispatch({ type: actions.TeethDataActions.ADD_DIAGNOSIS, payload: diagnosis });
                },
                err => console.log(err));
    }

    removeDiagnosis(diagnosisId: string, toothNo: string, patientId: string) : void {
        
        var diagnosis = <IToothActionData>{
            id: diagnosisId,
            patientId: patientId, 
            toothNo: toothNo, 
        }

        this.http.post('/api/toothStatus/removeToothDiagnosis', diagnosis).pipe(
            take(1))
            .subscribe(
                _=> {
                    this.store.dispatch({ type: actions.TeethDataActions.REMOVE_DIAGNOSIS, payload: diagnosisId });
                },
                err => console.log(err));
    }

    addManipulation(toothNo: string, manipulationId: number, patientId: string) : void {

        var manupulation = <IToothActionData>{
            id: GuidService.newGuid(),
            patientId:patientId, 
            toothNo: toothNo, 
            actionId: Number(manipulationId),
            date: new Date()
        }

        this.http.post('/api/toothStatus/addToothManipulation', manupulation).pipe(
            take(1))
            .subscribe(
                _=> {
                    this.store.dispatch({ type: actions.TeethDataActions.ADD_MANIPULATION, payload: manupulation });
                },
                err => console.log(err));
    }

    removeManupulation(manipulationId: string, toothNo: string, patientId: string) : void {
        var manipulation = <IToothActionData>{
            id: manipulationId,
            patientId: patientId, 
            toothNo: toothNo, 
        }

        this.http.post('/api/toothStatus/removeToothManipulation', manipulation).pipe(
            take(1))
            .subscribe(
                _=> {
                    this.store.dispatch({ type: actions.TeethDataActions.REMOVE_MANIPULATION, payload: manipulationId });
                },
                err => console.log(err));
    }
}