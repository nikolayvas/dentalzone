
import {take} from 'rxjs/operators';
 import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { FormGroup, ReactiveFormsModule, FormBuilder, Validators, FormControl } from "@angular/forms";

import { AuthService } from '../../auth/auth.service';
import { NotificationsManager} from "../../services/notifications-manager";
import { InputValidators } from '../../validation/input-validators'

@Component({
    selector: 'reset-password',
    templateUrl: './reset-password.component.html'
  })
  export class ResetPasswordComponent implements OnInit {

    form: FormGroup;
    email: FormControl;
    newPassword: FormControl;
    confirmPassword: FormControl;
    code: FormControl;

    formSubmitted: boolean = false;

    constructor(
        private router: Router,
        private authService: AuthService,
        private notificationsManager: NotificationsManager
       
    ) {}

    ngOnInit() {

        this.authService.logout();
    
        this.form = new FormGroup({
          email: this.email = new FormControl("", Validators.compose([Validators.required, InputValidators.validateEmail])),
          newPassword: this.newPassword = new FormControl("", Validators.required),
          confirmPassword: this.confirmPassword = new FormControl("", Validators.required),
          code: this.code = new FormControl("", Validators.required)
        });
      }

    reset() {
        if(this.newPassword.value !== this.confirmPassword.value) {
            this.notificationsManager.WarningNotification("Password does not match the confirm password.");
            return;
        }

        this.authService.resetPassword$(this.email.value, this.newPassword.value, this.code.value ).pipe(
            take(1))
            .subscribe(
                resposne => {
                    if(resposne){
                        this.notificationsManager.ServerSuccess("Your password was successfully reset.");
                        this.router.navigate(['app/login']);
                    } else {
                        this.notificationsManager.WarningNotification("Unsuccessful reset! Possible reasons could be non existing email or expired or not matching verification code.");
                    }
                },
                error => {
                    this.notificationsManager.ServerError(error, "reset-password");
                });
    }
}