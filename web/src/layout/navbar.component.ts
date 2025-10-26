import { Component } from '@angular/core';

@Component({
	selector: 'layout-navbar',
	template: `
		<div class="logo">Kubedash</div>
		<hr/>
		<div class="links">
			<core-link linkTo="/home">Home</core-link>
			<core-link linkTo="/home">Clusters</core-link>
			<core-link linkTo="/home">Pods</core-link>
			<core-link linkTo="/home">Services</core-link>
			<core-link linkTo="/home">Metrics</core-link>
		</div>
	`,
	styles: `
		:host {
			top: 0;
			left: 0;
			height: 100vh;
			border-right: 1px solid #e7e9ea;
			background: white;
			box-shadow: 0 2px 12px 0 rgba(11,22,44,0.05);
			padding: 10px 25px;
		}

		.logo {
			font-size: 2em;
			font-weight: 500;
		}

		.links {
			display: flex;
			flex-direction: column;
			justify-content: start;
		}

		.links core-link {
			margin-top: 20px;
		}
	`
})
export class LayoutNavbarComponent {}
