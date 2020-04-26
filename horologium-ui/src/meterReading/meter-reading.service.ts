import {Injectable} from '@angular/core';
import {MeterReadingServiceListener} from './meter-reading-service-listener';
import gql from 'graphql-tag';
import {map} from 'rxjs/operators';
import {Observable} from 'rxjs';
import {MeterReading} from './meter-reading';
import {Apollo} from 'apollo-angular';
import {Plan} from '../plan/plan';
import * as moment from 'moment';
import {MutationOptions} from 'apollo-client';

@Injectable({
  providedIn: 'root'
})
export class MeterReadingService {

  private listeners: MeterReadingServiceListener[] = [];

  constructor(private apollo: Apollo) {
  }

  public addListener(listener: MeterReadingServiceListener): void {
    this.listeners.push(listener);
  }

  queryMeterReadings(series: number): Observable<MeterReading[]> {
    return this.apollo.query({
      query: gql`
        query queryReadings($seriesId: Int!){
          meterReadings(seriesId: $seriesId){id, seriesId, count, date}
        }`,
      variables: {seriesId: series},
      errorPolicy: 'all',
      fetchPolicy: 'network-only'
    }).pipe(map(response => (response.data as any).meterReadings));
  }

  saveMeterReading(meterReading: MeterReading): Observable<MeterReading> {
    const from = moment(meterReading.date).format('YYYY-MM-DD');
    let mutation;
    if (meterReading.id) {
      mutation = this.getMutationForExistingReading(meterReading);
      return this.apollo.mutate(mutation).pipe(map(response => {
          const modifiedReading: MeterReading = (response.data as any).modifyMeterReading;
          if (modifiedReading) {
            for (const listener of this.listeners) {
              listener.meterReadingChanged(modifiedReading);
            }
          }
          return modifiedReading;
        }
      ));
    } else {
      mutation = this.getMutationForNewReading(meterReading);
      return this.apollo.mutate(mutation).pipe(map(response => {
          const createdReading: MeterReading = (response.data as any).createMeterReading;
          if (createdReading) {
            for (const listener of this.listeners) {
              listener.meterReadingAdded(createdReading);
            }
          }
          return createdReading;
        }
      ));
    }
  }

  private getMutationForNewReading(meterReading: MeterReading): MutationOptions {
    const from = moment(meterReading.date).format('YYYY-MM-DD');
    return {
      mutation: gql`
        mutation create($readingObj: MeterReadingInput!){
          createMeterReading(reading: $readingObj){id, seriesId, date, count}
        }`,
      variables: {
        readingObj: {
          count: meterReading.count,
          date: from,
          seriesId: meterReading.seriesId
        }
      },
      errorPolicy: 'all'
    };
  }

  private getMutationForExistingReading(meterReading: MeterReading): MutationOptions {
    const from = moment(meterReading.date).format('YYYY-MM-DD');
    return {
      mutation: gql`
        mutation modify($id: Int!, $readingObj: MeterReadingChange!){
          modifyMeterReading(id: $id, reading: $readingObj){id, seriesId, date, count}
        }`,
      variables: {
        readingObj: {
          count: meterReading.count,
          date: from,
        },
        id: meterReading.id
      },
      errorPolicy: 'all'
    };
  }
}
