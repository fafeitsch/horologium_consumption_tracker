import { Injectable } from '@angular/core';
import {Observable} from 'rxjs';
import {Series} from '../series/series';
import gql from 'graphql-tag';
import {map} from 'rxjs/operators';
import {Apollo} from 'apollo-angular';
import * as moment from 'moment';

@Injectable({
  providedIn: 'root'
})
export class StatisticsService {

  constructor(private apollo: Apollo) { }

  getMonthlyStatistics(seriesId: number, start: Date, end: Date): Observable<Statistics[]> {
    return this.apollo.query({
      query: gql`query stats($seriesId: Int!, $start: Date!, $end: Date!){
        monthlyStatistics(seriesId: $seriesId, start: $start, end: $end){
          validFrom, validTo, costs, consumption
        }
      }`,
      variables: {
        seriesId: seriesId,
        start: moment(start).format('YYYY-MM-DD'),
        end: moment(end).format('YYYY-MM-DD')
      },
    }).pipe(map(response => (response.data as any).monthlyStatistics));
  }
}
