import { Injectable } from '@angular/core';

@Injectable()
export class TokenStorageService {
    private _token: string;

    set token(value: string) {
        this._token = value;

        try {
            window.sessionStorage['token'] = value;
        }
        catch (e) {
            //do nothing
        }
    }

    get token(): string {

        let result: string = null;
        try {
            result = window.sessionStorage['token'];
        }
        catch (e) {
            result = this._token;
        }

        return result;
    }
}