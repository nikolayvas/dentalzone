export class AuthUserModel {
    constructor(private fEmail: string, private fName: string) {
    }

    get email(): string {
        return this.fEmail;
    }

    get name(): string {
        return this.fName;
    }
}