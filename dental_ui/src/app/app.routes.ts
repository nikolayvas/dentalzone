import { ModuleWithProviders } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { HomeComponent } from './components/home/home.component'
import { LoginComponent } from './components/login/login.component'
import { ForgotPasswordComponent } from './components/login/forgot-password.component'
import { ResetPasswordComponent } from './components/login/reset-password.component'
import { SignUpComponent } from './components/sign-up/signup.component'
import { SignUpActivateComponent } from './components/sign-up/signup-activate.component'
import { ClientPortalComponent} from './components/client-portal/client-portal.component'
import { PatientsListComponent } from './components/patient-list/patients-list.component'
import { PatientProfileComponent } from './components/patient-details/add-edit-patient.component'
import { ToothStatusComponent } from './components/tooth-status/tooth-status.component'
import { ScheduleComponent } from './components/schedule/schedule.component'

export const appRoutes: Routes = [
    { path: '', redirectTo: '/app', pathMatch: 'full' },
    { path: 'app', component: HomeComponent },
    { path: 'app/login', component: LoginComponent },
    { path: 'app/login/forgot-password', component: ForgotPasswordComponent},
    { path: 'app/login/reset-password', component: ResetPasswordComponent},
    { path: 'app/signup', component: SignUpComponent },
    { path: 'app/activate', component: SignUpActivateComponent },
    { path: 'app/portal', component: ClientPortalComponent, children: [
        { path: 'patients', component: PatientsListComponent},
        { path: 'schedule', component: ScheduleComponent},
        { path: 'patients/edit-patient-profile/:id', component: PatientProfileComponent},
        { path: 'patients/add-patient-profile', component: PatientProfileComponent},
        { path: 'patients/tooth-status/:id', component: ToothStatusComponent}
        ]
    }
]

export const appRouting: ModuleWithProviders = RouterModule.forRoot(appRoutes);