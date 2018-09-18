
import {map} from 'rxjs/operators';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { Observable } from 'rxjs';
import { AuthService } from '../../auth/auth.service';

@Component({
    templateUrl: 'client-portal.component.html'
})

export class ClientPortalComponent implements OnInit, OnDestroy {
    isLogged$: Observable<boolean>;
    userName$: Observable<string>;

    menuOpen: boolean = false;
    menuShort: boolean = false;

    constructor(private auth: AuthService
    ) {
        this.isLogged$ = auth.user$.pipe(map(n => !!n))
        this.userName$ = auth.user$.pipe(map(n => {if (n) {return n.name}}));
    }

    ngOnInit() {
   
    }

    ngOnDestroy() {
        
    }

    private logout() {
        this.auth.logout();
    }
 
    onActivate() {
        // a workaround to scroll to top when a navigation occurs
        //window.scrollTo(0, 0)
    }

    toggleOpenCloseMenu() {
        this.menuOpen = !this.menuOpen;
    }

    toggleShortWideMenu() {
        this.menuShort = !this.menuShort;
    }
}
