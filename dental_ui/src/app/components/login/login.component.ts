
import {take} from 'rxjs/operators';
import { Component, OnInit, ViewContainerRef } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { FormGroup, ReactiveFormsModule, FormBuilder, Validators, FormControl } from "@angular/forms";

import { AuthService } from '../../auth/auth.service';
import { NotificationsManager} from "../../services/notifications-manager";
import { InputValidators } from '../../validation/input-validators'

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html'
})
export class LoginComponent implements OnInit {

  form: FormGroup;
  email: FormControl;
  password: FormControl;

  formSubmitted: boolean = false;

  constructor(
    private router: Router,
    private authService: AuthService,
    private notificationsManager: NotificationsManager
  ) {
    
  }

  ngOnInit() {

    this.authService.logout();

    this.form = new FormGroup({
      email: this.email =  new FormControl("", Validators.compose([Validators.required, InputValidators.validateEmail])),
      password: this.password = new FormControl("", Validators.required)
    });
  }

  login() {

    this.formSubmitted = true;

    if(!this.form.valid) {
      return;
    }

    this.authService.login$(this.email.value, this.password.value).pipe(
      take(1))
      .subscribe(
        resposne => {
          if(resposne){
              this.router.navigate(['app/portal']);
          } else {
              this.notificationsManager.WarningNotification("Wrong username or password");
          }
        },
        error => {
            this.notificationsManager.ServerError(error, "login");
        });
  }
}
