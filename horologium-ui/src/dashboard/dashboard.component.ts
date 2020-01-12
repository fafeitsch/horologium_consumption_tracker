import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import {Series} from '../series/series';
import {SeriesService} from '../series/series.service';
import {PricingPlanService} from '../plan/pricing-plan.service';
import {Plan} from '../plan/plan';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {

  // noinspection JSMismatchedCollectionQueryUpdate
  private series: Series[];
  private selectedSeries: Series;
  private pricingPlans: Plan[];

  constructor(private seriesService: SeriesService, private planService: PricingPlanService, private router: Router) {
  }

  ngOnInit() {
    this.seriesService.getAllSeries().subscribe(resp => {
      this.series = resp;
    }, (error) => {
      this.router.navigate(['login']).then();
    });
  }

  private selectedSeriesChanged(series: Series): void {
    this.selectedSeries = series;
    this.planService.queryPricingPlans(series.id).subscribe(resp => {
      this.pricingPlans = resp;
      console.log(this.pricingPlans)
    }, (error) => {
      this.router.navigate(['login']).then();
    });
  }
}
