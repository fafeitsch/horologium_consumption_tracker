import {LoginService} from './login.service';
import {Component, OnInit} from '@angular/core';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent implements OnInit {

  constructor(private loginService: LoginService) {
  }

  ngOnInit() {
  }

  public attemptLogin(username: string, password: string) {
    this.loginService.startLogin(username, password);
  }

}
