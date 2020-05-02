import {NgModule} from '@angular/core';
import {ApolloModule, APOLLO_OPTIONS} from 'apollo-angular';
import {HttpLinkModule, HttpLink} from 'apollo-angular-link-http';
import {InMemoryCache} from 'apollo-cache-inmemory';
import {environment} from '../environments/environment';
import {ApolloLink} from 'apollo-link';
import {setContext} from 'apollo-link-context';
import {DefaultOptions} from 'apollo-client';

const uri = environment.graphQlServer; // <-- add the URL of the GraphQL server here

export function createApollo(httpLink: HttpLink) {
  const basic = setContext((operation, context) => ({
    headers: {
      Accept: 'charset=utf-8',
      'Content-Type': 'application/json'
    }
  }));

  const defaultOptions: DefaultOptions = {
    watchQuery: {
      fetchPolicy: 'no-cache',
      errorPolicy: 'ignore',
    },
    query: {
      fetchPolicy: 'no-cache',
      errorPolicy: 'all',
    },
  };

  // Get the authentication token from local storage if it exists
  const token = localStorage.getItem('token'); // 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyTmFtZSI6ImpvaG4ifQ.YAhFVqwtUicCpXp0IUYYx04Lv9Nt9t39UoufA4-Ic_E';
  const auth = setContext((operation, context) => ({
    headers: {
      Authorization: `Bearer ${token}`
    },
  }));

  const link = ApolloLink.from([basic, auth, httpLink.create({uri: uri})]);
  const cache = new InMemoryCache();

  return {
    link: link,
    cache: cache,
    defaultOptions: defaultOptions
  };
}


@NgModule({
  exports: [ApolloModule, HttpLinkModule],
  providers: [
    {
      provide: APOLLO_OPTIONS,
      useFactory: createApollo,
      deps: [HttpLink],
    },
  ],
})
export class GraphQLModule {
}
