import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Plan} from '../plan';
import {faPlus} from '@fortawesome/free-solid-svg-icons';
import {PricingPlanEditorComponent} from '../pricing-plan-editor/pricing-plan-editor.component';
import {Observable} from 'rxjs';

@Component({
  selector: 'app-pricing-plan-manager',
  templateUrl: './pricing-plan-manager.component.html',
  styleUrls: ['./pricing-plan-manager.component.scss']
})
export class PricingPlanManagerComponent {

  public plusIcon = faPlus;
  public selectedPlan: Plan;
  @Input() public pricingPlans: Plan[];
  @Input() public savePlan: (plan: Plan) => Observable<Plan>;

  constructor() {
  }

  onAddClicked(): void {
    // const dialogRef = this.dialog.open(PricingPlanEditorComponent, {
    //   height: '480px',
    //   width: '400px',
    //   data: {
    //     savePressed: this.savePlan
    //   },
    // });
  }
}
