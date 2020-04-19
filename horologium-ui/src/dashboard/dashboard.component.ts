import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import {Series} from '../series/series';
import {SeriesService} from '../series/series.service';
import {PricingPlanService} from '../plan/pricing-plan.service';
import {Plan} from '../plan/plan';
import {Observable} from 'rxjs';
import {PricingPlanServiceListener} from '../plan/pricing-plan-service-listener';
import {MeterReading} from '../meterReading/meter-reading';
import {MeterReadingServiceListener} from '../meterReading/meter-reading-service-listener';
import {MeterReadingService} from '../meterReading/meter-reading.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit, PricingPlanServiceListener, MeterReadingServiceListener {

  // noinspection JSMismatchedCollectionQueryUpdate
  public series: Series[];
  public selectedSeries: Series;
  public pricingPlans: Plan[];
  public meterReadings: MeterReading[];
  public savePlan: (plan: Plan) => Observable<Plan>;
  public saveMeterReading: (reading: MeterReading) => Observable<MeterReading>;

  constructor(private seriesService: SeriesService,
              private planService: PricingPlanService,
              private meterReadingService: MeterReadingService,
              private router: Router) {
    planService.addListener(this);
    meterReadingService.addListener(this);
    this.savePlan = (plan: Plan) => {
      plan.seriesId = this.selectedSeries.id;
      return this.planService.savePricingPlan(plan);
    };
    this.saveMeterReading = (reading: MeterReading) => {
      reading.seriesId = this.selectedSeries.id;
      return this.meterReadingService.saveMeterReading(reading);
    };
  }

  ngOnInit() {
    this.seriesService.getAllSeries().subscribe(resp => {
      this.series = resp;
    }, (error) => {
      console.log(error);
      this.router.navigate(['login']).then();
    });
  }

  public selectedSeriesChanged(series: Series): void {
    this.selectedSeries = series;
    this.planService.queryPricingPlans(series.id).subscribe(resp => {
      this.pricingPlans = resp;
    }, (error) => {
      console.log(error);
      this.router.navigate(['login']).then();
    });
    this.meterReadingService.queryMeterReadings(series.id).subscribe(resp => {
      this.meterReadings = resp;
    }, (error) => {
      console.log(error);
      this.router.navigate(['login']).then();
    });
  }

  public pricingPlanAdded(newPlan: Plan): void {
    if (newPlan && this.selectedSeries && newPlan.seriesId === this.selectedSeries.id) {
      // trigger manual update
      this.selectedSeriesChanged(this.selectedSeries);
    }
  }

  public meterReadingAdded(meterReading: MeterReading): void {
    if (meterReading && this.selectedSeries && meterReading.seriesId === this.selectedSeries.id) {
      // trigger manual update
      this.selectedSeriesChanged(this.selectedSeries);
    }
  }
}
