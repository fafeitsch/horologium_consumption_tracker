import {Injectable} from '@angular/core';
import {Apollo} from 'apollo-angular';
import {Observable} from 'rxjs';
import {Plan} from './plan';
import gql from 'graphql-tag';
import {map} from 'rxjs/operators';
import * as moment from 'moment';

@Injectable({
  providedIn: 'root'
})
export class PricingPlanService {

  constructor(private apollo: Apollo) {
  }

  queryPricingPlans(series: number): Observable<Plan[]> {
    return this.apollo.query({
      query: gql`
        query queryPlans($seriesId: Int!){
          pricingPlans(seriesId: $seriesId){name, basePrice, id, seriesId, unitPrice, validFrom, validTo}
        }`,
      variables: {seriesId: series},
      errorPolicy: 'all'
    }).pipe(map(response => (response.data as any).pricingPlans));
  }

  savePricingPlan(plan: Plan): Observable<Plan> {
    const from = moment(plan.validFrom).format('YYYY-MM-DD');
    let to = null;
    if (plan.validTo) {
      to = moment(plan.validTo).format('YYYY-MM-DD');
    }
    return this.apollo.mutate({
      mutation: gql`
        mutation create($planObj: NewPricingPlanInput){
          createPricingPlan(plan: $planObj){name, basePrice, id, unitPrice, validFrom, validTo}
        }`,
      variables: {
        planObj: {
          name: plan.name,
          seriesId: plan.seriesId,
          basePrice: plan.basePrice,
          unitPrice: plan.unitPrice,
          validFrom: from,
          validTo: to,
        }
      },
      errorPolicy: 'all'
    }).pipe(map(response => (response.data as any).createPricingPlan));
  }
}
