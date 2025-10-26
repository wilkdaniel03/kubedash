import { Component } from '@angular/core';

@Component({
	selector: 'layout-header',
	template: `
		<div class="wrapper">
			<div>Home</div>
			<core-input [rounded]="true" placeholder="Search"/>
		</div>
	`,
	styles: `
		:host {
			display: flex;
			justify-content: center;
			align-items: center;
			font-size: 1.5em;
			padding: 20px 0;
			border-bottom: 1px solid #e7e9ea;
			box-shadow: 0 2px 12px 0 rgba(11,22,44,0.05);
			background: #fff;
		}

		.wrapper {
			width: 70%;
			display: flex;
			justify-content: space-between;
			align-items: center;;
			column-gap: 30px;
		}
	`
})
export class LayoutHeaderComponent {}
