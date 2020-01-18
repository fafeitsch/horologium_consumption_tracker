import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Plan} from '../plan';
import {faPlus} from '@fortawesome/free-solid-svg-icons';
import {PricingPlanEditorComponent} from '../pricing-plan-editor/pricing-plan-editor.component';
import {MatDialog, MatDialogRef, MAT_DIALOG_DATA} from '@angular/material/dialog';
import {Observable} from 'rxjs';

@Component({
  selector: 'app-pricing-plan-manager',
  templateUrl: './pricing-plan-manager.component.html',
  styleUrls: ['./pricing-plan-manager.component.scss']
})
export class PricingPlanManagerComponent implements OnInit {

  private plusIcon = faPlus;
  @Input() private pricingPlans: Plan[];
  @Input() private savePlan: (plan: Plan) => Observable<Plan>;

  constructor(private dialog: MatDialog) {
  }

  ngOnInit() {
  }

  onAddClicked(): void {
    const dialogRef = this.dialog.open(PricingPlanEditorComponent, {
      height: '480px',
      width: '400px',
      data: {
        savePressed: this.savePlan
      },
    });
    dialogRef.beforeClosed().subscribe((result: boolean) => {
      console.log(result);
      console.log(dialogRef.componentInstance.getPlan());
    });
  }
}
