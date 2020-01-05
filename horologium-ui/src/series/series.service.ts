import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';
import {Series} from './series';
import gql from 'graphql-tag';
import {Apollo} from 'apollo-angular';
import {map} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class SeriesService {

  constructor(private apollo: Apollo) {
  }

  getAllSeries(): Observable<Series[]> {
    return this.apollo.query({
      query: gql`query {
        allSeries{id, name}
      }`,
      errorPolicy: 'all'
    }).pipe(map(response => (response.data as any).allSeries));
  }
}
