import {Component, OnInit} from '@angular/core';
import {LoginObserver} from '../login/loginObserver';
import {LoginService} from '../login/login.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit, LoginObserver {

  private loggedInAs: string;
  private loggedIn: boolean;

  constructor(private loginService: LoginService) {
    loginService.registerObserver(this);
  }

  ngOnInit(): void {
    const userName = this.loginService.getUsernameIfPresent();
    this.loggedIn = userName != null;
    if (userName === null) {
      this.loggedInAs = 'Not logged in';
    } else {
      this.loggedInAs = 'Logged in as ' + userName;
    }
  }

  successFullLogin(username: string): void {
    this.loggedInAs = 'Logged in as ' + username;
    this.loggedIn = true;
  }

  loggedOut(): void {
    this.ngOnInit();
  }
}
