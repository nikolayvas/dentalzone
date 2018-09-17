import { Pipe, PipeTransform } from '@angular/core';
import { IPatientData } from '../../models/patient.dto';

@Pipe({
    name: 'filter',
    pure: false
  })
export class FilterPatientsListPipe implements PipeTransform {
    transform(items: IPatientData[], query: string): IPatientData[] {
        if (!items || !query) {
            return items;
        }
        const results = items.filter(item => this.filter(item, query.toLowerCase()));
        
        return results;
    }

    filter(item: IPatientData, query: string): boolean {
        if (!!item) {
            return item.firstName.toLowerCase().includes(query) ||
            item.firstName.toLowerCase().includes(query) ||
            item.middleName.toLowerCase().includes(query) ||
            item.lastName.toLowerCase().includes(query) ||
            item.address.toLowerCase().includes(query) ||
            item.email.toLowerCase().includes(query) ||
            item.phoneNumber.toLowerCase().includes(query) || 
            item.generalInfo.toLowerCase().includes(query)
        }

        return false;
    }
}