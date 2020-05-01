import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppComponent} from './app.component';
import {GraphQLModule} from './graphql.module';
import {HttpClientModule} from '@angular/common/http';
import {HeaderComponent} from '../header/header.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {
  DateAdapter,
  MAT_DATE_FORMATS,
  MatButtonModule,
  MatCardModule, MatDatepickerModule, MatDialogModule,
  MatFormFieldModule,
  MatGridListModule,
  MatInputModule,
  MatListModule, MatRadioModule, MatSidenavModule, MatTableModule, MatTabsModule,
  MatToolbarModule
} from '@angular/material';
import {MatMomentDateModule} from '@angular/material-moment-adapter';
import {LoginComponent} from '../login/login.component';
import {DashboardComponent} from '../dashboard/dashboard.component';
import {SeriesListComponent} from '../series/list/series-list.component';
import {FontAwesomeModule} from '@fortawesome/angular-fontawesome';
import {PricingPlanTableComponent} from '../plan/pricing-plan-table/pricing-plan-table.component';
import {PricingPlanManagerComponent} from '../plan/pricing-plan-manager/pricing-plan-manager.component';
import {PricingPlanEditorComponent} from '../plan/pricing-plan-editor/pricing-plan-editor.component';
import {FormsModule} from '@angular/forms';
import {MY_DATE_FORMATS} from './MyDateAdapter';
import { MeterReadingManagerComponent } from '../meterReading/meter-reading-manager/meter-reading-manager.component';
import { MeterReadingTableComponent } from '../meterReading/meter-reading-table/meter-reading-table.component';
import { MeterReadingEditorComponent } from '../meterReading/meter-reading-editor/meter-reading-editor.component';
import { StatisticsComponentComponent } from '../statistics/statistics-component/statistics-component.component';
import { StatisticsTableComponent } from '../statistics/statistics-table/statistics-table.component';
import { InfoComponent } from '../info-component/info.component';

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
    MeterReadingManagerComponent,
    MeterReadingTableComponent,
    MeterReadingEditorComponent,
    StatisticsComponentComponent,
    StatisticsTableComponent,
    InfoComponent,
  ],
  imports: [
    BrowserModule,
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
    MatDialogModule,
    MatMomentDateModule,
    FormsModule,
  ],
  entryComponents: [
    PricingPlanEditorComponent,
    MeterReadingEditorComponent,
  ],
  providers: [
    {provide: MAT_DATE_FORMATS, useValue: MY_DATE_FORMATS},
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
