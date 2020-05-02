import {Component, Input, OnChanges, OnInit, SimpleChanges, ViewChild} from '@angular/core';
import {Series} from '../../series/series';
import {StatisticsService} from '../statistics.service';
import {DateAdapter, MAT_DATE_FORMATS, MAT_DATE_LOCALE, MatDatepicker} from '@angular/material';
import {Moment} from 'moment';
import {MAT_MOMENT_DATE_FORMATS, MomentDateAdapter} from '@angular/material-moment-adapter';

@Component({
  selector: 'app-statistics-component',
  templateUrl: './statistics-component.component.html',
  styleUrls: ['./statistics-component.component.scss']
})
export class StatisticsComponentComponent implements OnInit, OnChanges {
  @Input() public series: Series;
  public startDate: Date;
  public endDate: Date;
  public currentStats: Statistics[];

  constructor(private statisticService: StatisticsService) {
  }

  ngOnInit() {
    const year: number = new Date().getFullYear();
    this.startDate = new Date(year + '-01-01');
    const month: number = new Date().getMonth();
    this.endDate = new Date(year + '-' + month + '-01');
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (this.series) {
      this.load();
    }
  }

  public load(): void {
    this.statisticService.getMonthlyStatistics(this.series.id, this.startDate, this.endDate).subscribe((stats: Statistics[]) => this.currentStats = stats);
  }
}
