import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';
import {Plan} from '../../plan/plan';
import {MeterReading} from '../meter-reading';

@Component({
  selector: 'app-meter-reading-editor',
  templateUrl: './meter-reading-editor.component.html',
  styleUrls: ['./meter-reading-editor.component.scss']
})
export class MeterReadingEditorComponent implements OnInit {

  private date: Date;
  private count: number;
  private validityResult: boolean;

  constructor(
    private dialogRef: MatDialogRef<MeterReadingEditorComponent>,
    @Inject(MAT_DIALOG_DATA) private data: any) {
  }

  ngOnInit() {

  }

  public getMeterReading(): MeterReading {
    const reading: MeterReading = new MeterReading();
    reading.count = this.count;
    reading.date = this.date;
    return reading;
  }

  private saveClicked(): void {
    const reading: MeterReading = this.getMeterReading();
    this.data.savePressed(reading).subscribe((result: Plan) => {
      this.validityResult = null;
      this.dialogRef.close();
    }, error => {
      this.validityResult = error;
    });
  }
}
