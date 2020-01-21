export class MeterReading {
  private _count: number;
  private _date: Date;

  get count(): number {
    return this._count;
  }

  set count(value: number) {
    this._count = value;
  }

  get date(): Date {
    return this._date;
  }

  set date(value: Date) {
    this._date = value;
  }
}
