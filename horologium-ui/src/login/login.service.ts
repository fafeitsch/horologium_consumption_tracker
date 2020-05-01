import {Injectable} from '@angular/core';
import {HttpClient, HttpErrorResponse, HttpHeaders} from '@angular/common/http';
import {environment} from '../environments/environment';
import {LoginObserver} from './loginObserver';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class LoginService {
  private observers: LoginObserver[] = [];

  constructor(private http: HttpClient) {
  }

  public registerObserver(observer: LoginObserver): void {
    this.observers.push(observer);
  }

  public getUsernameIfPresent(): string {
    const token = localStorage.getItem('token');
    if (token === null) {
      return null;
    }
    return this.parseToken(token).UserName;
  }

  private parseToken(token: string): any {
    return JSON.parse(atob(token.split('.')[1]));
  }

  public logout(): void {
    localStorage.removeItem('token');
    this.observers.forEach(obs => obs.loggedOut());
  }

  public startLogin(username: string, password: string) {
    const headers = new HttpHeaders().set('Authorization', 'Basic ' + btoa(username + ':' + password));
    const httpResult = this.http.get<any>(environment.loginServer, {headers, observe: 'response'});
    httpResult.subscribe(resp => {
        localStorage.setItem('token', resp.body.token);
        this.observers.forEach(obs => obs.successFullLogin(this.parseToken(resp.body.token).UserName));
      },
      error => {
        const httpError = error as HttpErrorResponse;
        console.log(httpError.message);
      });
  }
}
