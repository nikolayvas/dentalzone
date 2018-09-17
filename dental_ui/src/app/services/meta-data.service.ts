import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { StoreService } from './store.service'
import { IMetaData } from '../models/meta-data.dto'
import { ManipulationModel } from '../models/manipulation.model';
import { DiagnosisModel } from '../models/diagnosis.mode';
import { ToothStatusModel } from '../models/tooth-status.model';
import * as actions from '../store/seed-actions';

@Injectable()
export class MetaDataService {

    constructor(
        private http: HttpClient,
        private store: StoreService
    ) {
        this.seedData();
    }

    seedData() {
        const result = this.http.get<IMetaData>('/api/seedData');

        result.subscribe(resp=> {
            if(resp) {
                this.store.dispatch({ type: actions.ClientPortalSeedActions.RESET_STATE_DATA, payload: resp });
            }
        })
    }

    manipulations$(): Observable<ManipulationModel[]> {
        return this.store.select(n => n.data.clientPortalStore.metaData.manipupations);
    }

    diagnosis$(): Observable<DiagnosisModel[]> {
        return this.store.select(n => n.data.clientPortalStore.metaData.diagnosis);
    }

    toothStatuses$(): Observable<ToothStatusModel[]> {
        return this.store.select(n => n.data.clientPortalStore.metaData.toothStatuses);
    }
}
