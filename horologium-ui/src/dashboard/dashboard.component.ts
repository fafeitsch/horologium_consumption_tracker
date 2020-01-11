import { Component, OnInit } from '@angular/core';
import {Router} from '@angular/router';
import {Series} from '../series/series';
import {SeriesService} from '../series/series.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {

  // noinspection JSMismatchedCollectionQueryUpdate
  private series: Series[];
  private selectedSeries: Series;

  constructor(private seriesService: SeriesService, private router: Router) { }

  ngOnInit() {
    this.seriesService.getAllSeries().subscribe(resp => {
      this.series = resp;
    }, (error) => {
      this.router.navigate(['login']).then();
    });
  }

  private selectedSeriesChanged(series: Series): void {
    this.selectedSeries = series;
    console.log(this.selectedSeries);
  }
}
