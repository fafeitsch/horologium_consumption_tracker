export class MeterReading {
  private _id: number;
  private _count: number;
  private _date: Date;
  private _seriesId: number;

  get id(): number {
    return this._id;
  }

  set id(value: number) {
    this._id = value;
  }

  get seriesId(): number {
    return this._seriesId;
  }

  set seriesId(value: number) {
    this._seriesId = value;
  }

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
