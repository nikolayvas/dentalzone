import { Component, Input } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';

@Component({
    selector: 'extended-control',
    templateUrl: 'extended-control.component.html'
})
export class ExtendedControlComponent {

    @Input() control: FormControl;
    @Input() formsubmitted: boolean;
    @Input() labeltext: string;
    @Input() classLabel: string;
    @Input() classValue: string;

    private get displayLabel(): boolean {
        return this.labeltext && this.labeltext.trim() !== '';
    }

    // the div in the template will only be added if
    // the control is dirty or form is submitted
    isDisplayed() {
        return (this.control.dirty || this.formsubmitted) && !this.control.valid && !this.control.disabled;
    }
}


