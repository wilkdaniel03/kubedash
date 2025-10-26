import { Directive, ViewContainerRef, OnInit, OnDestroy, inject } from '@angular/core';
import { CoreFormComponent } from '../../core';
import { Subscription, tap } from 'rxjs';
import { AuthenticationService, HttpService } from '@services';
import { FORM_FIELD_TYPE } from '@interfaces';
import { ApiError, AUTH_STATUS } from '@interfaces';
import { Router } from '@angular/router';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';

@Directive({
	selector: '[loginFormCfg]'
})
export class LoginFormCfg implements OnInit, OnDestroy {
	private _authenticationService = inject(AuthenticationService);
	private _router = inject(Router);
	private _httpService = inject(HttpService);
	private ref$: Subscription;
	private _containerRef: ViewContainerRef;

	constructor(
		public authenticationService: AuthenticationService,
		public router: Router,
		public http: HttpService,
		public httpService: HttpService,
		public containerRef: ViewContainerRef
	) {
		this._authenticationService = authenticationService;
		this._router = router;
		this._httpService = httpService;
		this._containerRef = containerRef;
	}

	async ngOnInit() {
		const ref = this._containerRef.createComponent(CoreFormComponent);
		ref.instance.model = [
			{ type: FORM_FIELD_TYPE.INPUT_TEXT, name: 'user', placeholder: 'username' },
			{ type: FORM_FIELD_TYPE.INPUT_PASSWORD, name: 'pass', placeholder: 'password' }
		]
		ref.instance.submitContent = "Sign Up";
		//ref.instance.submitContent = "Sign In!";
		this.ref$ = ref.instance.submit.subscribe(res => {
			const username: string = res['user'];
			const password: string = res['pass'];
			this._httpService.login({user:username,pass:password})
				.pipe(tap({
					next: res => {
						this._authenticationService.setStatus(AUTH_STATUS.OK);
						this._authenticationService.setUser({user:username,pass:password});
						this._authenticationService.setToken(res.token);
						this._router.navigate(['/home']);
						console.log(res);
					},
					error: (err: HttpErrorResponse) => {
						this._authenticationService.setStatus(AUTH_STATUS.FAILED);
						this._authenticationService.setUser({user:username,pass:password});
						console.log(err.error as ApiError);
					}
				})).subscribe();
		});
	}

	async ngOnDestroy() {
		this.ref$.unsubscribe();
	}
}
