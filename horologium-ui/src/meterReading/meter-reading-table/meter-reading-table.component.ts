import {Component, Input, OnInit} from '@angular/core';
import {Plan} from '../../plan/plan';
import {MeterReading} from '../meter-reading';
import {faPen} from '@fortawesome/free-solid-svg-icons';
import {Observable} from 'rxjs';
import {MeterReadingEditorComponent} from '../meter-reading-editor/meter-reading-editor.component';
import {MatDialog} from '@angular/material';

@Component({
  selector: 'app-meter-reading-table',
  templateUrl: './meter-reading-table.component.html',
  styleUrls: ['./meter-reading-table.component.scss']
})
export class MeterReadingTableComponent implements OnInit {

  public penIcon = faPen;
  @Input() public meterReadings: MeterReading;
  public displayedColumns: string[] = ['date', 'count', 'edit'];
  @Input() public saveMeterReading: (reading: MeterReading) => Observable<MeterReading>;

  constructor(private dialog: MatDialog) {
  }

  ngOnInit() {
  }

  public editMeterReading(reading: MeterReading): void {
    const dialogRef = this.dialog.open(MeterReadingEditorComponent, {
      height: '300px',
      width: '400px',
      data: {
        savePressed: this.saveMeterReading,
        existing: reading
      },
    });
  }

}
