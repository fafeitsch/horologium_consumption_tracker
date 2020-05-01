import {Component, OnInit} from '@angular/core';
import {LoginService} from '../login/login.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit{
  title = 'horologium-ui';

  public loggedIn: boolean;

  constructor(public loginService: LoginService) {
    loginService.registerObserver(this);
  }

  ngOnInit(): void {
    const userName = this.loginService.getUsernameIfPresent();
    this.loggedIn = userName != null;
  }

  successFullLogin(username: string): void {
    this.loggedIn = true;
  }

  loggedOut(): void {
    this.ngOnInit();
  }
}
