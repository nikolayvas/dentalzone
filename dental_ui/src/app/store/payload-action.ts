import { Action } from '@ngrx/store';

export class PayloadAction implements Action {
    type: string;

    constructor(public payload?: any) {
    }
}
