import {Component, Input, OnInit} from '@angular/core';
import {Plan} from '../../plan/plan';
import {MeterReading} from '../meter-reading';

@Component({
  selector: 'app-meter-reading-table',
  templateUrl: './meter-reading-table.component.html',
  styleUrls: ['./meter-reading-table.component.scss']
})
export class MeterReadingTableComponent implements OnInit {

  @Input() public meterReadings: MeterReading;
  public displayedColumns: string[] = ['date', 'count'];

  constructor() {
  }

  ngOnInit() {
  }

}
