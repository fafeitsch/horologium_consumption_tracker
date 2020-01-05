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

  constructor(private seriesService: SeriesService, private router: Router) { }

  ngOnInit() {
    this.seriesService.getAllSeries().subscribe(resp => {
      this.series = resp;
      console.log(this.series);
    }, (error) => {
      this.router.navigate(['login']).then();
    });
  }
}
