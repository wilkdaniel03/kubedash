import { Component, HostListener, inject, Input } from '@angular/core';
import { Router } from '@angular/router';

@Component({
	selector: 'core-link',
	template: `
		<div><ng-content></ng-content></div><i class="fa-solid fa-chevron-right"></i>
	`,
	styles: `
		:host {
			font-size: 1.4em;
			color: black;
			cursor: pointer;
			display: flex;
			justify-content: space-between;
			align-items: center;
		}

		.fa-chevron-right {
			font-size: 1em;
		}

		:host:hover,:host:focus {
			color: #0c5dec;
		}
	`
})
export class CoreLinkComponent {
	private _router = inject(Router);

	@Input()
	linkTo: string = '';

	@HostListener('click')
	onClick(): void {
		this._router.navigate([this.linkTo]);
	}
}
