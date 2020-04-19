import {Component, Input, OnInit} from '@angular/core';
import {Series} from '../../series/series';
import {StatisticsService} from '../statistics.service';

@Component({
  selector: 'app-statistics-component',
  templateUrl: './statistics-component.component.html',
  styleUrls: ['./statistics-component.component.scss']
})
export class StatisticsComponentComponent implements OnInit {

  @Input() public series: Series;
  public startDate: Date;
  public endDate: Date;
  public currentStats: Statistics[];

  constructor(private statisticService: StatisticsService) {
  }

  ngOnInit() {
  }

  public load(): void {
    this.statisticService.getMonthlyStatistics(this.series.id, this.startDate, this.endDate).subscribe((stats: Statistics[]) => this.currentStats = stats);
  }

}
