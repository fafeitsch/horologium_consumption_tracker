import {Injectable} from '@angular/core';
import {Apollo} from 'apollo-angular';
import {Observable} from 'rxjs';
import {Plan} from './plan';
import gql from 'graphql-tag';
import {map} from 'rxjs/operators';

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
}
