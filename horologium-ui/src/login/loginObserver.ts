export interface LoginObserver {
  successFullLogin(username: string): void;
  loggedOut(): void;
}
