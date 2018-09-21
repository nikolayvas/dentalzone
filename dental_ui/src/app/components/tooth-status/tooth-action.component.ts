import { Component, EventEmitter, Input, Output, OnInit } from '@angular/core';
import { IActionTypeModel } from '../../models/action-type.model'
import { ToothActionModel} from '../../models/tooth-action.model'
import { Utils } from '../../services/utils'
import { ConfirmationService } from 'primeng/primeng';

@Component({
    selector: 'tooth-action',
    templateUrl: './tooth-action.component.html'
})
export class ToothActionComponent implements OnInit {

    @Input()
    actionName: string;

    @Input()
    currentToothNo: string;

    @Input()
    toothActions: ToothActionModel[] = []

    @Input()
    actionsList: IActionTypeModel[] = []

    @Output()
    onAdd = new EventEmitter<{actionTypeId: number}>();

    @Output()
    onRemove = new EventEmitter<{toothActionid: string}>();

    private selectedAction: IActionTypeModel

    constructor(
        private confirmationService: ConfirmationService
    ) {
           
    }
        
    ngOnInit() {
    }

    getActionName(id: number): string {
        const item = this.actionsList.find(n => n.id == id)
        return item ? item.name : "";
    }

    addAction() {
        this.onAdd.emit({actionTypeId: this.selectedAction.id});
    }

    removeAction(id: string): void {
        this.confirmationService.confirm({
            header: 'Remove action',
            message: 'Are you sure you want to remove that item?',
            icon: 'fa fa-question-circle',
            accept: () => {
                this.onRemove.emit({toothActionid: id});
            }
        });
    }

    get allowAddAction(): boolean {
        return this.selectedAction && !Utils.isBlankOrEmpty(this.currentToothNo);
    }
}