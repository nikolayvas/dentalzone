
import {throwError as observableThrowError, of as observableOf} from 'rxjs';

import {catchError, map} from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpResponse, HttpErrorResponse } from '@angular/common/http';
import { Observable, BehaviorSubject} from 'rxjs';
import { AuthUserModel } from './auth-user.model';
import { TokenStorageService } from './token-storage.service';

@Injectable()
export class AuthService {

    private _user$ = new BehaviorSubject<AuthUserModel>(null);

    get user$(): Observable<AuthUserModel> {
        return this._user$;
    }

    constructor(
        private http: HttpClient,
        private tokenStorage: TokenStorageService
    ) {
        this.checkLoggedUser();
    }

    checkLoggedUser() {
        const result = this.http.get<{email: string, name: string}>('/api/dentist/get');

        result.subscribe(resp=> {
            if(resp.email) {
                this._user$.next(new AuthUserModel(resp.email, resp.name));
            }
        })
    }

    signup$(name: string, email: string, password: string): Observable<boolean> {
        const result = this.http.post<boolean>('/api/dentist/signup', JSON.stringify({ name, email, password }));
        
        return result.pipe(
            map(resp=> {
                return true;
            }),
            catchError((err: HttpErrorResponse) => {
                console.log(`An error occured: ${err.message}`)

                if(err.status == 406) {
                    return observableOf(false);
                }

                return observableThrowError(err);
            }),);
    }

    signupActivate$(id: string): Observable<boolean> {
        const result = this.http.get<boolean>('/api/dentist/activate', {params: {id: id}});
        
        return result.pipe(
            map(resp=> {
                return true;
            }),
            catchError((err: HttpErrorResponse) => {
                console.log(`An error occured: ${err.message}`)

                if(err.status == 406) {
                    return observableOf(false);
                }

                return observableThrowError(err);
            }),);
    }

    login$(email: string, password: string): Observable<boolean> {
        const result = this.http.post<{ email: string, name: string, token: string }>('/api/dentist/login', JSON.stringify({ email, password }));

        return result.pipe(
            map(resp=> {
                this.tokenStorage.token = resp.token;

                this._user$.next(new AuthUserModel(resp.email, resp.name));
                return true;
            }),
            catchError((err: HttpErrorResponse) => {
                console.log(`An error occured: ${err.message}`)

                //no such user error
                if(err.status == 404) {
                    return observableOf(false);
                }

                return observableThrowError(err);
            }),);
    }

    requestResetPasswordCode$(email: string): Observable<boolean> {
        const result = this.http.post<boolean>('/api/dentist/sendPasswordResetConfirmationCode', JSON.stringify({ email }));

        return result.pipe(
            map(resp=> {
                return true;
            }),
            catchError((err: HttpErrorResponse) => {
                console.log(`An error occured: ${err.message}`)

                if(err.status == 409) {
                    return observableOf(false);
                }

                return observableThrowError(err);
            }),);
    }

    resetPassword$(email: string, newPassword: string, confirmationCode: string): Observable<boolean>  {
        const result = this.http.post<boolean>('/api/dentist/resetPassword', JSON.stringify({ email, newPassword, confirmationCode }));

        return result.pipe(
            map(resp=> {
                return true;
            }),
            catchError((err: HttpErrorResponse) => {
                console.log(`An error occured: ${err.message}`)

                if(err.status == 409) {
                    return observableOf(false);
                }

                return observableThrowError(err);
            }),);
    }

    logout() {
        this.tokenStorage.token = "";
    }
}