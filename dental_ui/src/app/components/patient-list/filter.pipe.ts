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

            const parts = query.split(" ");

            for (var i = 0; i < parts.length; i++) {
               if(!this.hasMatch(item, parts[i] )) {
                   return false;
               }
            }

            return true;
        }

        return false;
    }

    hasMatch(item: IPatientData, query: string): boolean {
        if (!!item) {

            return (item.id && item.id.toLowerCase().includes(query)) ||
            (item.firstName && item.firstName.toLowerCase().includes(query)) ||
            (item.firstName && item.firstName.toLowerCase().includes(query)) ||
            (item.middleName && item.middleName.toLowerCase().includes(query)) ||
            (item.lastName && item.lastName.toLowerCase().includes(query)) ||
            (item.address && item.address.toLowerCase().includes(query)) ||
            (item.email && item.email.toLowerCase().includes(query)) ||
            (item.phoneNumber && item.phoneNumber.toLowerCase().includes(query)) || 
            (item.generalInfo && item.generalInfo.toLowerCase().includes(query))
        }
    }
}