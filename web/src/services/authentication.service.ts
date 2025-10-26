import { Injectable } from '@angular/core';
import { IUser, AUTH_STATUS } from '@interfaces';

@Injectable()
export class AuthenticationService {
	private user: IUser = {
		user: 'Daniel',
		pass: ''
	};
	private status: AUTH_STATUS = AUTH_STATUS.OK;
	private token: string = "";

	public setUser(user:IUser): void {
		this.user.user = user.user;
		this.user.pass = user.user;
	}

	public setStatus(status: AUTH_STATUS): void {
		this.status = status;
	}

	public setToken(token: string): void {
		this.token = token;
	}

	public getUser(): IUser {
		return this.user;
	}

	public getStatus(): AUTH_STATUS{
		return this.status;
	}

	public getToken(): string {
		return this.token;
	}
}
