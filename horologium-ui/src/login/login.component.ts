import {LoginService} from './login.service';
import {Router} from '@angular/router';
import {Component, OnInit} from '@angular/core';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
})
export class LoginComponent implements OnInit {

  constructor(private loginService: LoginService, private router: Router) {
  }

  ngOnInit() {
  }

  private attemptLogin(username: string, password: string) {
    this.loginService.startLogin(username, password);
  }

}
