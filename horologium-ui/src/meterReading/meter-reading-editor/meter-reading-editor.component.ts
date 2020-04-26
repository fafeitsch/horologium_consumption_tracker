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

  public date: Date;
  public count: number;
  public validityResult: boolean;

  constructor(
    private dialogRef: MatDialogRef<MeterReadingEditorComponent>,
    @Inject(MAT_DIALOG_DATA) private data: any) {
  }

  ngOnInit() {
    console.log(this.data)
    if (this.data.existing) {
      this.count = this.data.existing.count;
      this.date = this.data.existing.date;
    }
  }

  private getMeterReading(): MeterReading {
    let reading: MeterReading;
    if (this.data.existing) {
      reading = this.data.existing;
      reading.count = this.count;
      reading.date = this.date;
      return reading;
    } else {
      return {
        id: undefined,
        count: this.count,
        date: this.date,
        seriesId: undefined,
      };
    }
  }

  public saveClicked(): void {
    const reading: MeterReading = this.getMeterReading();
    this.data.savePressed(reading).subscribe((result: MeterReading) => {
      this.validityResult = null;
      this.dialogRef.close();
    }, error => {
      this.validityResult = error;
    });
  }
}
