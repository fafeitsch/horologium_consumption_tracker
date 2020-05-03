export interface Plan {
  id: number;
  seriesId: number;
  name: string;
  basePrice: number;
  unitPrice: number;
  validFrom: Date;
  validTo: Date;
}
