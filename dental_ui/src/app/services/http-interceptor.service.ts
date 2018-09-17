
import {throwError as observableThrowError} from 'rxjs';

import {catchError} from 'rxjs/operators';
import { Injectable, Injector } from '@angular/core';
import { HttpRequest, HttpHandler, HttpInterceptor, HttpEvent, HttpErrorResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment'
import { Router } from '@angular/router';

import { TokenStorageService } from '../auth/token-storage.service'
import { NotificationsManager } from '../services/notifications-manager'


@Injectable()
export class CustomHttpInterceptor implements HttpInterceptor {

    constructor(private injector: Injector) {
        
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

        return next.handle(req).pipe(
            catchError((response: any)=> {
                if (response instanceof HttpErrorResponse) {
                    const notification = this.injector.get(NotificationsManager);

                    console.log('response in the catch: ', response);
                    notification.ServerError(response);

                    if(response.status == 401) {
                        this.injector.get(Router).navigate(['app/login']);
                    }
                }
        
                return observableThrowError(response);
        }));
  }
}