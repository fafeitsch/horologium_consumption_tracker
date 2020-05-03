import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';
import {Component, EventEmitter, Inject, Input, OnChanges, OnInit, Output, SimpleChange, SimpleChanges} from '@angular/core';
import {Plan} from '../plan';
import {Observable} from 'rxjs';

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

  constructor() {
  }

  ngOnChanges(changes: SimpleChanges): void {
    const planChange: SimpleChange = changes.pricingPlan;
    this.pricingPlan = planChange.currentValue;
    this.editedPlan = JSON.parse(JSON.stringify(planChange.currentValue));
  }

  public saveClicked(): void {
    this.savePressed(this.editedPlan).subscribe((result: Plan) => {
      this.validityResult = null;
    }, error => {
      this.validityResult = error;
    });
  }

  public resetClicked(): void {
    this.editedPlan = JSON.parse(JSON.stringify(this.pricingPlan));
  }
}


