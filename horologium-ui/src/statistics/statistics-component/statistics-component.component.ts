import {Component, Input, OnInit} from '@angular/core';
import {Series} from '../../series/series';
import {StatisticsService} from '../statistics.service';

@Component({
  selector: 'app-statistics-component',
  templateUrl: './statistics-component.component.html',
  styleUrls: ['./statistics-component.component.scss']
})
export class StatisticsComponentComponent implements OnInit {

  @Input() private series: Series;
  private startDate: Date;
  private endDate: Date;
  private currentStats: Statistics[];

  constructor(private statisticService: StatisticsService) {
  }

  ngOnInit() {
  }

  private load(): void {
    this.statisticService.getMonthlyStatistics(this.series.id, this.startDate, this.endDate).subscribe((stats: Statistics[]) => this.currentStats = stats);
  }

}
