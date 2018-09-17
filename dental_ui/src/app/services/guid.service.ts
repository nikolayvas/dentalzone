const emptyGuid = '00000000-0000-0000-0000-000000000000';

export class GuidService {
    static get empty(): string {
        return emptyGuid;
    }

    static isEmpty(guid: string): boolean {
        return guid === emptyGuid;
    }

    static newGuid(): string {
        //
        // code below has been copied from somewhere
        // not tested
        //
        var d = new Date().getTime();
        if (window.performance && typeof window.performance.now === "function") {
            d += performance.now(); //use high-precision timer if available
        }
        var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
            var r = (d + Math.random() * 16) % 16 | 0;
            d = Math.floor(d / 16);
            return (c === 'x' ? r : (r & 0x3 | 0x8)).toString(16);
        });
        return uuid;
    }
}