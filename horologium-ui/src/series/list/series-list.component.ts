import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Series} from '../series';
import {faPlay, faHeartbeat} from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-series-list',
  templateUrl: './series-list.component.html',
  styleUrls: ['./series-list.component.scss']
})
export class SeriesListComponent implements OnInit {

  private faHeartbeat = faHeartbeat;
  @Input() public series: Series[];
  @Output() public selectionChanged = new EventEmitter<Series>();
  constructor() { }

  ngOnInit() {
  }

  public seriesSelectionChanged(series: Series): void {
    this.selectionChanged.emit(series);
  }
}
