import {Plan} from './plan';

export interface PricingPlanServiceListener {
    pricingPlanAdded(newPlan: Plan): void;
}
