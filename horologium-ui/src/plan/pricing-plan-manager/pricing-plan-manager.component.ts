import {Component, Input, OnInit} from '@angular/core';
import {Plan} from '../plan';
import { faPlus } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-pricing-plan-manager',
  templateUrl: './pricing-plan-manager.component.html',
  styleUrls: ['./pricing-plan-manager.component.scss']
})
export class PricingPlanManagerComponent implements OnInit {

  private plusIcon = faPlus;
  @Input() private pricingPlans: Plan[];

  constructor() {
  }

  ngOnInit() {
  }

}
