import { Subscription } from 'rxjs';

export class Utils {
    static noop = () => { };
    
    static isBlank(obj: any): boolean {
        return obj === undefined || obj === null;
    }

    static isBlankOrEmpty(obj: any): boolean {
        return Utils.isBlank(obj) || obj === '' || (Array.isArray(obj) && obj.length === 0);
    }

    static isBlankOrWhiteSpace(obj: any): boolean {
        return Utils.isBlankOrEmpty(obj) || obj.toString().trim() === '';
    }

    static unsubscribe(subscriptions: Subscription | Subscription[]) {
        if (subscriptions) {
            if (Array.isArray(subscriptions)) {
                subscriptions.forEach(x => {
                    if (x && !x.closed) {
                        x.unsubscribe();
                    }
                });
            }
            else {
                subscriptions.unsubscribe();
            }
        }
    }
}