import {Injectable} from '@angular/core';
import {Apollo} from 'apollo-angular';
import {Observable} from 'rxjs';
import {Plan} from './plan';
import gql from 'graphql-tag';
import {map} from 'rxjs/operators';
import * as moment from 'moment';
import {PricingPlanServiceListener} from './pricing-plan-service-listener';
import {MutationOptions} from 'apollo-client';
import {MeterReading} from '../meterReading/meter-reading';

@Injectable({
  providedIn: 'root'
})
export class PricingPlanService {

  private listeners: PricingPlanServiceListener[] = [];

  constructor(private apollo: Apollo) {
  }

  addListener(listener: PricingPlanServiceListener): void {
    this.listeners.push(listener);
  }

  queryPricingPlans(series: number): Observable<Plan[]> {
    return this.apollo.query({
      query: gql`
        query queryPlans($seriesId: Int!){
          pricingPlans(seriesId: $seriesId){name, basePrice, id, seriesId, unitPrice, validFrom, validTo, seriesId}
        }`,
      variables: {seriesId: series},
      errorPolicy: 'all',
      fetchPolicy: 'network-only'
    }).pipe(map(response => (response.data as any).pricingPlans));
  }

  savePricingPlan(plan: Plan): Observable<Plan> {
    let mutation;
    if (plan.id) {
      mutation = this.getMutationForExistingPricingPlan(plan);
      return this.apollo.mutate(mutation).pipe(map(response => {
        return (response.data as any).modifyPricingPlan;
      }));
    } else {
      mutation = this.getMutationForNewPricingPlan(plan);
      return this.apollo.mutate(mutation).pipe(map(response => {
        const addedPlan: Plan = (response.data as any).modifyPricingPlan;
        if (addedPlan) {
          for (const listener of this.listeners) {
            listener.pricingPlanAdded(addedPlan);
          }
        }
        return addedPlan;
      }));
    }
  }

  private getMutationForExistingPricingPlan(plan: Plan): MutationOptions {
    const from = moment(plan.validFrom).format('YYYY-MM-DD');
    let to = null;
    if (plan.validTo) {
      to = moment(plan.validTo).format('YYYY-MM-DD');
    }
    return {
      mutation: gql`
        mutation modify($planObj:  PricingPlanChange!){
          modifyPricingPlan(plan: $planObj){id, name, validFrom, validTo, basePrice, unitPrice, seriesId}
        }`,
      variables: {
        planObj: {
          id: plan.id,
          name: plan.name,
          validFrom: from,
          validTo: to,
          basePrice: plan.basePrice,
          unitPrice: plan.unitPrice
        }
      },
      errorPolicy: 'all'
    };
  }

  private getMutationForNewPricingPlan(plan: Plan): MutationOptions {
    const from = moment(plan.validFrom).format('YYYY-MM-DD');
    let to = null;
    if (plan.validTo) {
      to = moment(plan.validTo).format('YYYY-MM-DD');
    }
    return {
      mutation: gql`
        mutation create($planObj:  PricingPlanInput!){
          createPricingPlan(plan: $planObj){id, name, validFrom, validTo, basePrice, unitPrice, seriesId}
        }`,
      variables: {
        planObj: {
          name: plan.name,
          validFrom: from,
          validTo: to,
          basePrice: plan.basePrice,
          unitPrice: plan.unitPrice
        }
      },
      errorPolicy: 'all'
    };
  }
}
