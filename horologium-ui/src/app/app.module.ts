import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {GraphQLModule} from './graphql.module';
import {HttpClientModule} from '@angular/common/http';
import {HeaderComponent} from '../header/header.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {MatButtonModule, MatCardModule, MatFormFieldModule, MatGridListModule, MatInputModule, MatToolbarModule} from '@angular/material';
import {LoginComponent} from '../login/login.component';
import {DashboardComponent} from '../dashboard/dashboard.component';

@NgModule({
  declarations: [
    AppComponent,
    HeaderComponent,
    LoginComponent,
    DashboardComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    GraphQLModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MatInputModule,
    MatToolbarModule,
    MatCardModule,
    MatFormFieldModule,
    MatGridListModule,
    MatButtonModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {
}
