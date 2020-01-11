export class Plan {
  private _id: number;
  private _name: string;
  private _basePrice: number;
  private _unitPrice: number;
  private _validFrom: Date;
  private _validTo: Date;


  get id(): number {
    return this._id;
  }

  set id(value: number) {
    this._id = value;
  }

  get name(): string {
    return this._name;
  }

  set name(value: string) {
    this._name = value;
  }

  get basePrice(): number {
    return this._basePrice;
  }

  set basePrice(value: number) {
    this._basePrice = value;
  }

  get unitPrice(): number {
    return this._unitPrice;
  }

  set unitPrice(value: number) {
    this._unitPrice = value;
  }

  get validFrom(): Date {
    return this._validFrom;
  }

  set validFrom(value: Date) {
    this._validFrom = value;
  }

  get validTo(): Date {
    return this._validTo;
  }

  set validTo(value: Date) {
    this._validTo = value;
  }
}
