import {Component, Input, OnInit, Output} from '@angular/core';
import {Plan} from '../plan';
import {SelectionModel} from '@angular/cdk/collections';
import {EventEmitter} from '@angular/core';

@Component({
  selector: 'app-pricing-plan-table',
  templateUrl: './pricing-plan-table.component.html',
  styleUrls: ['./pricing-plan-table.component.scss']
})
export class PricingPlanTableComponent implements OnInit {

  @Input() public pricingPlans: Plan[];
  @Output() public selectedPlan: EventEmitter<Plan> = new EventEmitter();
  public selection: SelectionModel<Plan> = new SelectionModel<Plan>(false, []);
  public displayedColumns: string[] = ['select', 'name', 'basePrice', 'unitPrice', 'validFrom', 'validTo'];

  constructor() {
  }

  ngOnInit() {
  }

  public planClicked(plan: Plan): void {
    this.selection.toggle(plan);
    if (this.selection.isSelected(plan)) {
      this.selectedPlan.emit(plan);
    }
    if (this.selection.isEmpty()) {
      this.selectedPlan.emit(undefined);
    }
  }
}
