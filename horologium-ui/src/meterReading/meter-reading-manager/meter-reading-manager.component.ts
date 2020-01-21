import {Component, Input, OnInit} from '@angular/core';
import {MeterReading} from '../meter-reading';
import {faPlus} from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-meter-reading-manager',
  templateUrl: './meter-reading-manager.component.html',
  styleUrls: ['./meter-reading-manager.component.scss']
})
export class MeterReadingManagerComponent implements OnInit {

  private plusIcon = faPlus;
  @Input() meterReadings: MeterReading[];

  constructor() {
  }

  ngOnInit() {
  }

  onAddClicked(): void {

  }
}
