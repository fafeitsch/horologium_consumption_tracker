import {Component, Input, OnInit} from '@angular/core';
import {MeterReading} from '../meter-reading';
import {faPlus} from '@fortawesome/free-solid-svg-icons';
import {PricingPlanEditorComponent} from '../../plan/pricing-plan-editor/pricing-plan-editor.component';
import {MatDialog} from '@angular/material';
import {MeterReadingEditorComponent} from '../meter-reading-editor/meter-reading-editor.component';
import {Plan} from '../../plan/plan';
import {Observable} from 'rxjs';

@Component({
  selector: 'app-meter-reading-manager',
  templateUrl: './meter-reading-manager.component.html',
  styleUrls: ['./meter-reading-manager.component.scss']
})
export class MeterReadingManagerComponent implements OnInit {

  public plusIcon = faPlus;
  @Input() public meterReadings: MeterReading[];
  @Input() public saveMeterReading: (reading: MeterReading) => Observable<MeterReading>;

  constructor(private dialog: MatDialog) {
  }

  ngOnInit() {
  }

  onAddClicked(): void {
    const dialogRef = this.dialog.open(MeterReadingEditorComponent, {
      height: '300px',
      width: '400px',
      data: {
        savePressed: this.saveMeterReading
      },
    });
  }
}
