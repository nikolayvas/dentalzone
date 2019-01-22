
import {throwError as observableThrowError,  Observable } from 'rxjs';

import {catchError, finalize} from 'rxjs/operators';
import { Injectable, Injector } from '@angular/core';
import { HttpRequest, HttpHandler, HttpInterceptor, HttpEvent, HttpErrorResponse } from '@angular/common/http';
import { environment } from '../../environments/environment'
import { Router } from '@angular/router';

import { TokenStorageService } from '../auth/token-storage.service'
import { NotificationsManager } from '../services/notifications-manager'
import { ProgressIndicatorService } from './progress-indicator.service';


@Injectable()
export class CustomHttpInterceptor implements HttpInterceptor {

    constructor(
        private injector: Injector,
        private progress: ProgressIndicatorService
        ) {
        
    }

    intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
        const url = environment.serverHost;
        
        const tokenStorage = this.injector.get(TokenStorageService);
        
        req = req.clone({
            url: url + req.url
        });

        if(req.url.endsWith("/api/dentist/login") || req.url.endsWith("/api/dentist/signup") ){ 
            //TODO make another middleware for login where authentication is not required

            //we skipp global error handling for login/signup
            return next.handle(req)
        }

        req = req.clone({
            setHeaders: {
                Authorization: `Bearer ${tokenStorage.token}`
            }
        });

        if(req.method == "GET") {
            this.progress.set({isActive: true, infinite: true});
        }

        return next.handle(req).pipe(
            catchError((response: any)=> {
                if (response instanceof HttpErrorResponse) {

                    const notification = this.injector.get(NotificationsManager);
                    if(response.status == 401) {
                        notification.WarningNotification("Your session has been expired! You need to login again.");
                        this.injector.get(Router).navigate(['app/login']);
                    } 
                    else {
                        console.log('response in the catch: ', response);
                        notification.ServerError(response);
                    }
                }
        
                return observableThrowError(response);
        }),
        finalize(()=> {
            if(req.method == "GET") {
                this.progress.set({isActive: false});
            }
        }));
  }
}