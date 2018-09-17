
import {take} from 'rxjs/operators';
import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { FormGroup, ReactiveFormsModule, FormBuilder, Validators, FormControl } from "@angular/forms";

import { AuthService } from '../../auth/auth.service';
import { NotificationsManager} from "../../services/notifications-manager";
import { InputValidators } from '../../validation/input-validators'

@Component({
    selector: 'forgot-password',
    templateUrl: './forgot-password.component.html'
})
export class ForgotPasswordComponent implements OnInit {

    form: FormGroup;
    email: FormControl;

    formSubmitted: boolean = false;

    constructor(
        private router: Router,
        private authService: AuthService,
        private notificationsManager: NotificationsManager
       
    ) {}

    ngOnInit() {

        this.authService.logout();
    
        this.form = new FormGroup({
          email: this.email =  new FormControl("", Validators.compose([Validators.required, InputValidators.validateEmail])),
        });
      }

    search() {
        if(!this.form.valid) {
            return;
        }

        this.authService.requestResetPasswordCode$(this.email.value).pipe(
            take(1))
            .subscribe(
                resposne => {
                    if(resposne){
                        this.notificationsManager.ServerSuccess("Confirmation reset Code was successfully sent to your email.");
                        this.router.navigate(['app/login/reset-password']);
                    } else {
                        this.notificationsManager.WarningNotification("We couldn't find your account with that information");
                    }
                },
                error => {
                    this.notificationsManager.ServerError(error, "forgot-password");
                });
    }
}