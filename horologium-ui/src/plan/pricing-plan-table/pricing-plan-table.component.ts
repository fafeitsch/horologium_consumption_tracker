import {Component, Input, OnInit} from '@angular/core';
import {Plan} from '../plan';

@Component({
  selector: 'app-pricing-plan-table',
  templateUrl: './pricing-plan-table.component.html',
  styleUrls: ['./pricing-plan-table.component.scss']
})
export class PricingPlanTableComponent implements OnInit {

  @Input() public pricingPlans: Plan[];
  public displayedColumns: string[] = ['name', 'basePrice', 'unitPrice', 'validFrom', 'validTo'];

  constructor() {
  }

  ngOnInit() {
  }

}
