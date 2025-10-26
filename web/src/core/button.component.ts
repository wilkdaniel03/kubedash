import { Component } from '@angular/core';

@Component({
	selector: 'core-btn',
	template: `<button class="btn btn-outline-primary"><ng-content></ng-content></button>`,
	styles: `
		button {
			width: 100%;
			padding: 10px;
			font-weight: 700;
		}
	`
})
export class CoreButtonComponent {}
