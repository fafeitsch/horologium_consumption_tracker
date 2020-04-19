import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-statistics-table',
  templateUrl: './statistics-table.component.html',
  styleUrls: ['./statistics-table.component.scss']
})
export class StatisticsTableComponent implements OnInit {

  @Input() statistics: Statistics[];
  private displayedColumns: string[] = ['validFrom', 'validTo', 'consumption', 'costs'];

  constructor() { }

  ngOnInit() {
  }

}
