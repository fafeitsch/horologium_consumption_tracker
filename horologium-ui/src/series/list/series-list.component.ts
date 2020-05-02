import {Component, EventEmitter, Input, OnChanges, OnInit, Output} from '@angular/core';
import {Series} from '../series';
import {faPlay, faHeartbeat} from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-series-list',
  templateUrl: './series-list.component.html',
  styleUrls: ['./series-list.component.scss']
})
export class SeriesListComponent implements OnInit, OnChanges {

  private faHeartbeat = faHeartbeat;
  @Input() public series: Series[];
  @Output() public selectionChanged = new EventEmitter<Series>();

  constructor() {
  }

  ngOnInit() {
  }

  ngOnChanges() {
    if (this.series !== undefined && this.series.length > 0) {
      this.seriesSelectionChanged(this.series[0]);
    }
  }

  public seriesSelectionChanged(series: Series): void {
    this.selectionChanged.emit(series);
  }
}
