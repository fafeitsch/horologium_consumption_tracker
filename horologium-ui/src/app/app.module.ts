import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {GraphQLModule} from './graphql.module';
import {HttpClientModule} from '@angular/common/http';
import {HeaderComponent} from '../header/header.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {
  MatButtonModule,
  MatCardModule, MatDatepickerModule,
  MatFormFieldModule,
  MatGridListModule,
  MatInputModule,
  MatListModule, MatNativeDateModule, MatRadioModule, MatSidenavModule, MatTableModule, MatTabsModule,
  MatToolbarModule
} from '@angular/material';
import {LoginComponent} from '../login/login.component';
import {DashboardComponent} from '../dashboard/dashboard.component';
import { SeriesListComponent } from '../series/list/series-list.component';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { PricingPlanTableComponent } from '../plan/pricing-plan-table/pricing-plan-table.component';
import { PricingPlanManagerComponent } from '../plan/pricing-plan-manager/pricing-plan-manager.component';
import { PricingPlanEditorComponent } from '../plan/pricing-plan-editor/pricing-plan-editor.component';

@NgModule({
  declarations: [
    AppComponent,
    HeaderComponent,
    LoginComponent,
    DashboardComponent,
    SeriesListComponent,
    PricingPlanTableComponent,
    PricingPlanManagerComponent,
    PricingPlanEditorComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    GraphQLModule,
    HttpClientModule,
    BrowserAnimationsModule,
    FontAwesomeModule,
    MatInputModule,
    MatToolbarModule,
    MatCardModule,
    MatFormFieldModule,
    MatGridListModule,
    MatButtonModule,
    MatListModule,
    MatRadioModule,
    MatTableModule,
    MatSidenavModule,
    MatTabsModule,
    MatDatepickerModule,
    MatNativeDateModule,
  ],
  entryComponents: [
    PricingPlanEditorComponent
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {
}
