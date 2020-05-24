import {Component, Input, OnChanges, SimpleChange, SimpleChanges} from '@angular/core';
import {Plan} from '../plan';
import {Observable} from 'rxjs';
import {MatSnackBar} from '@angular/material/snack-bar';

@Component({
  selector: 'app-pricing-plan-editor',
  templateUrl: './pricing-plan-editor.component.html',
  styleUrls: ['./pricing-plan-editor.component.scss']
})
export class PricingPlanEditorComponent implements OnChanges {

  @Input() pricingPlan: Plan;
  @Input() savePressed: (plan: Plan) => Observable<Plan>;
  public editedPlan: Plan;
  public validityResult: string;

  constructor(private snackBar: MatSnackBar) {
  }

  ngOnChanges(changes: SimpleChanges): void {
    const planChange: SimpleChange = changes.pricingPlan;
    this.pricingPlan = planChange.currentValue;
    this.editedPlan = JSON.parse(JSON.stringify(planChange.currentValue));
  }

  public saveClicked(): void {
    this.pricingPlan.name = this.editedPlan.name;
    this.pricingPlan.validFrom = this.editedPlan.validFrom;
    this.pricingPlan.validTo = this.editedPlan.validTo;
    this.pricingPlan.unitPrice = this.editedPlan.unitPrice;
    this.pricingPlan.basePrice = this.editedPlan.basePrice;
    this.savePressed(this.editedPlan).subscribe((result: Plan) => {
      this.validityResult = null;
      this.snackBar.open('Plan saved successfully.', undefined, {duration: 1000});
    }, error => {
      this.snackBar.open('An error occured: ' + error, undefined, {duration: 1000});
      this.validityResult = error;
    });
  }

  public resetClicked(): void {
    this.editedPlan = JSON.parse(JSON.stringify(this.pricingPlan));
  }
}


