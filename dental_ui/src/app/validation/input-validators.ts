import { AbstractControl, FormGroup, FormArray, Validators, ValidatorFn } from '@angular/forms';

import { Utils } from '../services/utils';

export interface IInputValidatorResult {
    [key: string]: { valid: boolean; }
}

export class InputValidators {

    private static regexPhone = /^(?:\([0-9]{3}\)|[0-9]{3})(?:[-.]?|\s*)(?:[0-9]{3})(?:[-.]?|\s*)(?:[0-9]{4})$/;
    private static regexEmail = /^[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/i;
    private static monthYearDateFormat = 'MM/YYYY';
    private static regexWebsite = /^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)?[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/)?$/i;

    static getValidatorErrorMessage(validatorName: string, validatorValue?: any) {

        if (validatorName === 'minlength' && Utils.isBlank(validatorValue))
            throw new Error('Validator value is required');

        // ReSharper disable once QualifiedExpressionMaybeNull
        const config = {
            'required': 'Required',
            'invalidPhone': 'Invalid Phone Number',
            'duplicatePhone': 'Duplicate phone number',
            'invalidEmail': 'Invalid Email',
            'invalidEmailOrPhone': 'Invalid Email or Phone Number',
            'invalidDate': 'Invalid Date',
            'invalidWebsite': 'Invalid web site',
            'invalidZipCode': 'Please enter a valid zip code',
            'duplicateEmail': 'Duplicate Email',
            'minlength': `Minimum length ${validatorValue.requiredLength}`,
            'maxlength': `Maximum length ${validatorValue.requiredLength}`,
            'invalidNumber': 'Invalid Number',
            'patternString1': 'Allowed input: Alpha . - \' ^ ` ',
            'patternString2': 'Allowed input: Alpha, Numeric . , - \' #',
            'patternString3': 'Allowed input: Alpha, Numeric, Special Character',
        };

        let result = config[validatorName];
        if (!result) {
            result = validatorValue && validatorValue.message || null;
        }

        return result;
    }

    static validatePhone(control: AbstractControl): IInputValidatorResult {

        if (Utils.isBlankOrWhiteSpace(control.value) || InputValidators.regexPhone.test(control.value)) {
            return null;
        }

        return { 'invalidPhone': { valid: false } };
    }

    static validatePhoneDuplicate(controls: AbstractControl[]): ValidatorFn {

        return (control: AbstractControl): IInputValidatorResult => {
            if (controls.filter(x => !!x && !!x.value && x.value === control.value).length > 1) {
                return { 'duplicatePhone': { valid: false } };
            } else {
                return null;
            }
        }
    }

    static validateEmail(control: AbstractControl): IInputValidatorResult {

        if (Utils.isBlankOrWhiteSpace(control.value) || InputValidators.regexEmail.test(control.value)) {
            return null;
        }

        return { 'invalidEmail': { valid: false } };
    }

    static validateWebsite(control: AbstractControl): IInputValidatorResult {
        
        if (Utils.isBlankOrWhiteSpace(control.value) || InputValidators.regexWebsite.test(control.value)) {
            return null;
        }

        return { 'invalidWebsite': { valid: false } };
    }

    static validateEmailOrPhone(control: AbstractControl): IInputValidatorResult {

        if (Utils.isBlankOrWhiteSpace(control.value) ||
            InputValidators.regexPhone.test(control.value) ||
            InputValidators.regexEmail.test(control.value)) {
            return null;
        }

        return { 'invalidEmailOrPhone': { valid: false } };
    }

    static validateDate(control: AbstractControl): IInputValidatorResult {

        const regexDate = /^\d{1,2}\/\d{1,2}\/(\d{2}|\d{4})$/;
        const date = Date.parse(control.value);

        if (Utils.isBlankOrWhiteSpace(control.value) || (regexDate.test(control.value) && !isNaN(date))) {
            return null;
        }

        return { 'invalidDate': { valid: false } };
    }
}
