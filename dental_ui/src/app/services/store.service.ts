import { Injectable } from '@angular/core';
import { Action, Store } from '@ngrx/store';

import { ClientPortalState } from '../store/root-reducer';

@Injectable()
export class StoreService extends Store<ClientPortalState> {

}
