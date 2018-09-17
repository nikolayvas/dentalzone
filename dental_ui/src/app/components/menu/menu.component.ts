
import {of as observableOf,  Observable } from 'rxjs';
import { Component, OnInit } from '@angular/core';

import { MenuItem } from 'primeng/primeng';

@Component({
    selector: 'app-menu',
    templateUrl: 'menu.component.html'
})
export class MenuComponent implements OnInit {
    private isLogged$: Observable<boolean>;
    private active: string;

    items: MenuItem[];

    constructor(
        //private authService: AuthService
    ) {
        this.isLogged$ = observableOf(true); //trueauthService.authenticated$;
    }

    ngOnInit() {
        this.items = [
            { label: 'Patients', icon: 'icon-dashboard', routerLink: 'patients', command: (event) => {} /*this.disableAllBut(event.item)*/ },
        ];
    }

    disableAllBut(except: any) {
        for (let item of this.items) {
            if (item.label !== except.label) {
                item.expanded = false;
            }
        }
    }
}