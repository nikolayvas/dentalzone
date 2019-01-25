import { Component, EventEmitter, HostListener, OnInit, OnDestroy, Output} from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';

import { Subscription } from 'rxjs';
import { distinctUntilChanged } from 'rxjs/operators';

import { Utils } from '../../services/utils'
import { IPaginatorModel, DayOrWeekMode } from './schedule.models';

import * as moment from 'moment';
import * as R from 'ramda';
import { SelectItem } from 'primeng/api';

@Component({
    selector: 'paginator',
    templateUrl: './paginator.component.html',
  })
export class PaginatorComponent implements OnInit, OnDestroy {

    private _current: moment.Moment;
    private _subscriptions: Subscription = new Subscription();
    private _smallSize: number = 650;

    @Output()
    onPeriodChanged = new EventEmitter<IPaginatorModel>();

    @HostListener('window:resize', ['$event'])
    onResize(event) {
        if(window.innerWidth < this._smallSize) {
            this.mode.patchValue(DayOrWeekMode.Day)
            this.bigScreen = false;
        } 
        else {
            this.bigScreen = true;
        };
    }

    selectedPeriod: string;
    dayAndWeek: SelectItem[];
    bigScreen: boolean = true;

    form: FormGroup;
    mode: FormControl;

    constructor(
    ) {
        this.form = new FormGroup({
            mode: this.mode = new FormControl(''),
        });

        this._subscriptions.add(this.mode.valueChanges.pipe(distinctUntilChanged(R.equals)).subscribe(value=> {
            this._current = moment();
            this.updateSelectedPeriod();
        }));
    }

    ngOnInit() {
        this.dayAndWeek = [{label:'Day', value: DayOrWeekMode.Day}, {label:'Week', value:DayOrWeekMode.Week}];
        
        this._current = moment();
        if (window.innerWidth > this._smallSize) {
            this.mode.patchValue(DayOrWeekMode.Week);
        }
        else {
            this.mode.patchValue(DayOrWeekMode.Day);
        }
    }

    ngOnDestroy() {
        Utils.unsubscribe(this._subscriptions);
    }

    prev() {
        if(this.mode.value == DayOrWeekMode.Week) {
            this._current.add(-7, 'days');
        }
        else {
            this._current.add(-1, 'days');
        }

        this.updateSelectedPeriod();
    }

    today() {
        this._current = moment();

        this.updateSelectedPeriod();
    }

    next() {
        if(this.mode.value == DayOrWeekMode.Week) {
            this._current.add(7, 'days');
        } 
        else {
            this._current.add(1, 'days');
        }

        this.updateSelectedPeriod();
    }

    private updateSelectedPeriod() {
        if(this.mode.value == DayOrWeekMode.Week) {

            const monday = this._current.clone().weekday(1);
            const sunday = this._current.clone().weekday(7);

            const month1 = monday.format("MMM");
            const day1 = monday.format("DD");
            const year1 = monday.format("YYYY");

            const month2 = sunday.format("MMM");
            const day2 = sunday.format("DD");
            const year2 = sunday.format("YYYY");
            
            if(year1 == year2) {
                this.selectedPeriod = `${month1} ${day1} - ${month2} ${day2}, ${year1}`;
            }
            else {
                this.selectedPeriod = `${month1} ${day1}, ${year1} - ${month2} ${day2}, ${year2}`;
            }
        }
        else {
            this.selectedPeriod = this._current.format("dddd, MMM DD, YYYY");
        }

        this.onPeriodChanged.next(<IPaginatorModel>{pageMode: this.mode.value, currentDate: this._current})
    }
}