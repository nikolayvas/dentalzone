import { Component, EventEmitter, Input, Output, OnInit } from '@angular/core';
import { IActionTypeModel } from '../../models/action-type.model'
import { ToothActionModel} from '../../models/tooth-action.model'
import { Utils } from '../../services/utils'

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
        this.onRemove.emit({toothActionid: id});
    }

    get allowAddAction(): boolean {
        return this.selectedAction && !Utils.isBlankOrEmpty(this.currentToothNo);
    }
}