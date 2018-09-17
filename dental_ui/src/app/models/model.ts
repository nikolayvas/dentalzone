import { Iterable } from 'immutable';

export class KeyedModel<TDto extends Iterable<string, any>> {
    constructor(protected dto: TDto, addKeys: string[] = []) {

        let self = this;
        let proto = Object.getPrototypeOf(this);

        let obj = this.dto.toObject();
        let keys = [...Object.keys(obj), ...(addKeys || [])];

        keys.forEach((key: string) => {
            let descriptor = Object.getOwnPropertyDescriptor(proto, key);
            
            Object.defineProperty(self, key, {
                get: () => descriptor.get.call(self),
                enumerable: true,
                configurable: true
            });
        });
    }
}