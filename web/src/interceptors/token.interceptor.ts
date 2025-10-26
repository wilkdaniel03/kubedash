import { Injectable, inject } from '@angular/core';
import { HttpRequest, HttpInterceptor, HttpHandler } from '@angular/common/http';
import { AuthenticationService } from '@services';
import { BEARER_TOKEN_ENABLED } from './index';

@Injectable()
export class TokenInterceptor implements HttpInterceptor {
	private _authenticationService = inject(AuthenticationService);

	intercept(req: HttpRequest<any>, next: HttpHandler) {
		if(req.context.get(BEARER_TOKEN_ENABLED)) {
			const newReq = req.clone({
				setHeaders: {
					Authorization: `Bearer ${this._authenticationService.getToken()}`
				}
			});
			return next.handle(newReq);
		}
		return next.handle(req);
	}
}
