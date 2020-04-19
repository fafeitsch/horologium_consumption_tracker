import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';
import {Component, EventEmitter, Inject, OnInit, Output} from '@angular/core';
import {Plan} from '../plan';

@Component({
  selector: 'app-pricing-plan-editor',
  templateUrl: './pricing-plan-editor.component.html',
  styleUrls: ['./pricing-plan-editor.component.scss']
})
export class PricingPlanEditorComponent implements OnInit {

  public planName: string;
  public basePrice: number;
  public unitPrice: number;
  public validFrom: Date;
  public validTo: Date;
  public validityResult: string;

  constructor(
    public dialogRef: MatDialogRef<PricingPlanEditorComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any) {
  }

  ngOnInit() {

  }

  public getPlan(): Plan {
    const plan = new Plan();
    plan.name = this.planName;
    plan.basePrice = this.basePrice;
    plan.unitPrice = this.unitPrice;
    plan.validFrom = this.validFrom;
    plan.validTo = this.validTo;
    return plan;
  }

  public saveClicked(): void {
    const plan: Plan = this.getPlan();
    this.data.savePressed(plan).subscribe((result: Plan) => {
      this.validityResult = null;
      this.dialogRef.close();
    }, error => {
      this.validityResult = error;
    });
  }
}


