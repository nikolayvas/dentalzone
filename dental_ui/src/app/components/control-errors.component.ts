import { Component, Input } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';
import { InputValidators } from './../validation/input-validators';

@Component({
    selector: 'control-errors',
    template: `<span class="label label-danger m-t-xs inline-error-message" *ngIf="errorMessage !== null">{{errorMessage}}</span>`
})
export class ControlErrorsComponent {
    @Input() control: FormControl;
    @Input() formsubmitted: boolean;
    constructor() { }

    get errorMessage() {
        for (const propertyName in this.control.errors) {
            if (this.control.errors.hasOwnProperty(propertyName) && (this.control.dirty || this.formsubmitted)) {
                return InputValidators.getValidatorErrorMessage(propertyName, this.control.errors[propertyName]);
            }
        }

        return null;
    }
}
