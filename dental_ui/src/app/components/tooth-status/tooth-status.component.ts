
import {take} from 'rxjs/operators';
import { Component, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { Router, ActivatedRoute } from '@angular/router';
import { MetaDataService } from '../../services/meta-data.service'

import { ManipulationModel } from '../../models/manipulation.model';
import { DiagnosisModel } from '../../models/diagnosis.mode';
import { Utils } from '../../services/utils'

import { ToothStatusService } from '../../services/tooth-status.service'
import { ToothActionModel } from '../../models/tooth-action.model';

export interface IToothStatus {
    toothNo: number,
    status: string
}

@Component({
  selector: 'tooth-status',
  templateUrl: './tooth-status.component.html'
})
export class ToothStatusComponent implements OnInit {

    private subscriptions: Subscription = new Subscription();

    private patientId: string;
    private currentToothNo: string;

    private teeth: number[][] = [
        [18, 17, 16, 15, 14, 13, 12, 11, 0, 21, 22, 23, 24, 25, 26, 27, 28],
        [48, 47, 46, 45 ,44, 43, 42, 41, 0, 31, 32, 33, 34, 35, 36, 37, 38]
    ]

    private toothStatuses: IToothStatus[][] = [];
    private toothMap: { [tooth: number] : IToothStatus} = [];

    private manipulationsList: ManipulationModel[];
    private diagnosisList: DiagnosisModel[];

    private toothManipulations: ToothActionModel[] = [];
    private toothDiagnosis: ToothActionModel[] = []

    private selectedManipulation: ManipulationModel;
    private selectedDiagnosis: DiagnosisModel;

    get allowAddDiagnosis(): boolean {
        return this.selectedDiagnosis && !Utils.isBlankOrEmpty(this.currentToothNo);
    }

    get allowAddManipulation(): boolean {
        return this.selectedManipulation && !Utils.isBlankOrEmpty(this.currentToothNo);
    }

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        private metaDataService: MetaDataService,
        private toothService: ToothStatusService
      ) {
            this.patientId = this.route.snapshot.params['id'];

            for(var i: number = 0; i < this.teeth.length; i++) {
                this.toothStatuses[i] = [];
                for(var j: number = 0; j< this.teeth[i].length; j++) {
                    this.toothStatuses[i][j] = <IToothStatus>{toothNo: this.teeth[i][j]};

                    this.toothMap[this.teeth[i][j]] = this.toothStatuses[i][j];
                }
            }

            toothService.seedTeethData(this.patientId);
    }
        
    ngOnInit() {

        const manipulationsSubscription: Subscription = this.metaDataService.manipulations$().subscribe(n => {
            this.manipulationsList = n;
        })

        const diagnosisSubscription: Subscription = this.metaDataService.diagnosis$().subscribe(n => {
            this.diagnosisList = n;
        })

        const toothDiagnosisSubscription: Subscription = this.toothService.toothDiagnosisData$().subscribe(n=> {
            if(!this.currentToothNo) {
                this.toothDiagnosis = n
            }
            else {
                this.toothDiagnosis = n.filter(n=>n.toothNo == this.currentToothNo)
            }

            n.forEach(n=>this.toothMap[Number(n.toothNo)].status = "1")
        });

        const toothManipulationSubscription: Subscription = this.toothService.toothManipulationsData$().subscribe(n=> {
            if(!this.currentToothNo) {
                this.toothManipulations = n
            }
            else {
                this.toothManipulations = n.filter(n=>n.toothNo == this.currentToothNo)
            }

            n.forEach(n=>this.toothMap[Number(n.toothNo)].status = "1")
        });

        this.subscriptions.add(manipulationsSubscription);
        this.subscriptions.add(diagnosisSubscription);
        this.subscriptions.add(toothDiagnosisSubscription);
        this.subscriptions.add(toothManipulationSubscription);
    }

    ngOnDestroy() {
        Utils.unsubscribe(this.subscriptions);
        this.toothService.clearTeethData();
    }

    getToothData(toothNo?: string): void {
        this.toothService.toothDiagnosisData$(toothNo).pipe(take(1)).subscribe(n=>{
            this.toothDiagnosis = n
        });

        this.toothService.toothManipulationsData$(toothNo).pipe(take(1)).subscribe(n=>{
            this.toothManipulations = n
        })
    }

    toothClick(toothNo: string): void {
        this.currentToothNo = toothNo;
        this.getToothData(toothNo);
    }

    addDiagnosis(data: {actionTypeId: number}) {
        this.toothService.addDiagnosis(String(this.currentToothNo), data.actionTypeId, this.patientId);
    }

    removeDiagnosis(data: {toothActionid: string}): void {
        this.toothService.removeDiagnosis(data.toothActionid);
    }

    addManipulation(data: {actionTypeId: number}) {
        this.toothService.addManipulation(String(this.currentToothNo), data.actionTypeId, this.patientId);
    }

    removeManipulation(data: {toothActionid: string}): void {
        this.toothService.removeManupulation(data.toothActionid);
    }

    back() {
        this.router.navigateByUrl('/app/portal/patients');  
    }
} 