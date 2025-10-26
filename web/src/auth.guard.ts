import { inject, Injectable } from '@angular/core';
import { AuthenticationService } from '@services';
import { AUTH_STATUS } from '@interfaces';
import { CanActivate, Router } from '@angular/router';

@Injectable({ providedIn: "root" })
export class AuthGuard implements CanActivate {
	private _authenticationService = inject(AuthenticationService);
	private _router = inject(Router);

	canActivate(): boolean {
		const result = this._authenticationService.getStatus() === AUTH_STATUS.OK
			? true
			: false;
		if(!result)
			this._router.navigate(['/home']);
		return result;
	}
}
