import {Component, Input, OnInit} from '@angular/core';
import {Series} from '../series';
import {faPlay, faHeartbeat} from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-series-list',
  templateUrl: './series-list.component.html',
  styleUrls: ['./series-list.component.scss']
})
export class SeriesListComponent implements OnInit {

  private faHeartbeat = faHeartbeat;
  @Input() private series: Series[];
  constructor() { }

  ngOnInit() {
  }

  private seriesSelectionChanged(series: Series): void {
  }
}
