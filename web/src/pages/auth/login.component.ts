import { Component } from '@angular/core';

@Component({
	selector: 'auth-page',
	template: `
		<div class="box">
			<h1>Get Access</h1>
			<div><ng-container loginFormCfg></ng-container></div>
		</div>
	`,
	styles: `
		.box {
			width: 30%;
			height: 40%;
			position: absolute;
			top: 25%;
			left: 50%;
			transform: translateX(-50%);
		}

		.box div {
			height: 80%;
		}
	`
})
export class AuthPageComponent {}
