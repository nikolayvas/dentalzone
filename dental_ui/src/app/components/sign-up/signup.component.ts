
import {take} from 'rxjs/operators';
import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { FormGroup, ReactiveFormsModule, FormBuilder, Validators, FormControl } from "@angular/forms";

import { AuthService } from '../../auth/auth.service';
import { NotificationsManager} from "../../services/notifications-manager";
import { InputValidators } from '../../validation/input-validators'

@Component({
  selector: 'sign-upp',
  templateUrl: './signup.component.html'
})
export class SignUpComponent implements OnInit {

  form: FormGroup;
  name: FormControl;
  email: FormControl;
  password: FormControl;

  formSubmitted: boolean=false;

  constructor(
    private router: Router,
    private authService: AuthService,
    private notificationsManager: NotificationsManager
  ) {

  }

  ngOnInit() {
    this.authService.logout();

    this.form = new FormGroup({
        name: this.name =  new FormControl("", Validators.required),
        email: this.email =  new FormControl("", Validators.compose([Validators.required, InputValidators.validateEmail])),
        password: this.password = new FormControl("", Validators.required)
      });
  }

  signup() {

    if(!this.form.valid) {
        return;
    }

    this.authService.signup$(this.name.value, this.email.value, this.password.value).pipe(
      take(1))
      .subscribe(
        resposne => {
          if(resposne){
              this.router.navigate(['app/login']);
              this.notificationsManager.ServerSuccess("You have successfully signed up. Check your email please and follow the verification prosess.");
          } else {
              this.notificationsManager.WarningNotification("There is already registered user with such email !");
          }
        },
        error => {
            this.notificationsManager.ServerError(error, "signup");
        });
  }
}
