
import {take} from 'rxjs/operators';
import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import { AuthService } from '../../auth/auth.service';
import { NotificationsManager} from "../../services/notifications-manager";


@Component({
  selector: 'sign-up-activate',
  templateUrl: './signup-activate.component.html'

})
export class SignUpActivateComponent implements OnInit {

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private authService: AuthService,
    private notificationsManager: NotificationsManager
  ) {
  }

  ngOnInit() {
    this.authService.logout();

    this.route.queryParams.pipe(
        take(1))
        .subscribe(params => {
            var id = params['id'];

            if(id) {
                this.authService.signupActivate$(id).pipe(
                take(1))
                .subscribe(resposne => {
                    if(resposne){
                        this.notificationsManager.ServerSuccess("You have successfully activated your user.");
                    } else {
                        this.notificationsManager.WarningNotification("There is not such activation code or expiration time has been reached!");
                    }
                },
                error => {
                    //the error will be handled from interceptor
                    //this.notificationsManager.ServerError(error, "activate user");
                });
            }

            this.router.navigate(['app/login']); 
        });
  }
}
