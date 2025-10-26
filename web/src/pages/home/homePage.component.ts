import { Component, ElementRef, OnInit, TemplateRef, ViewChild, ViewContainerRef, inject } from '@angular/core';
import { Pod, Service } from '@interfaces';
import { AuthenticationService,KubeService } from '@services';
import { ResourceModalComponent } from '../../features/resourceModal/resourceModal.component';

@Component({
	template: `
		<resource-modal/>
		<layout-navbar/>
		<div class="wrapper">
			<layout-header/>
			<div class="grid-container">
				<div class="first"><h1>Hello, {{ username }}!</h1></div>
				<div class="second"><core-btn (click)="onClick()">Create resource</core-btn></div>
				<div class="third">
					<div>
						<i class="fa-solid fa-box"></i>
						<h2>Pods</h2>
					</div>
					@for(pod of pods; track pod) {
						<li><core-dot/>{{ pod.name }}</li>
					}
				</div>
				<div class="fourth">
					<div>
						<i class="fa-solid fa-up-right-and-down-left-from-center"></i>
						<h2>Services</h2>
					</div>
					@for(service of services; track service) {
						<li><core-dot/>{{ service.name }}</li>
					}
				</div>
			</div>
		</div>
	`,
	styles: `
		:host {
			display: flex;
			background: #f7f9fc;
		}

		layout-navbar {
			flex-grow: 1;
		}

		.wrapper {
			flex-grow: 6;
			display: flex;
			flex-direction: column;
		}

		.wrapper h1 {
			display: block;
		}

		.grid-container {
			width: 70%;
			flex-grow: 1;
			margin-left: auto;
			margin-right: auto;
			display: grid;
			grid-gap: 50px;
			grid-template-columns: repeat(4,1fr);
			grid-template-rows: repeat(6,200px);
		}

		.first {
			grid-column: 1;
			grid-row: 1;
			display: flex;
			flex-direction: column;
			justify-content: center;
		}

		.second {
			grid-column: 4;
			grid-row: 1;
			display: flex;
			flex-direction: column;
			justify-content: center;
		}

		.third {
			grid-row: 2 / span 2;
			grid-column: 1 / span 2;
			background: #fff;
			box-shadow: 1px 2px 12px 0 rgba(11,22,44,0.05);
			border-radius: 25px;
			padding: 20px 10px;
			border: 1px solid #e7e9ea;
			overflow: scroll;
		}

		.third i {
			color: #0c5dec;
			font-size: 2em;
			background: #e6e8eb;
			padding: 20px;
			border-radius: 10px;
		}

		.third > div {
			width: 50%;
			display: flex;
			justify-content: space-evenly;
			margin-bottom: 25px;
		}

		.fourth {
			grid-row: 2 / span 2;
			grid-column: 3 / span 2;
			background: #fff;
			box-shadow: 1px 2px 12px 0 rgba(11,22,44,0.05);
			border-radius: 25px;
			padding: 20px 10px;
			border: 1px solid #e7e9ea;
		}

		.fourth i {
			color: #0c5dec;
			font-size: 2em;
			background: #e6e8eb;
			padding: 20px;
			border-radius: 10px;
		}

		.fourth > div {
			width: 50%;
			display: flex;
			justify-content: space-evenly;
			margin-bottom: 25px;
		}

		li {
			list-style-type: none;
			padding: 10px 25px;
			display: flex;
			justify-content: start;
			align-items: center;
			font-size: 1.1em;
			border-radius: 15px;
			cursor: pointer;
		}

		li:hover {
			background: lightgray;
		}

		li core-dot {
			margin-right: 8px;
		}
	`
})
export class HomePageComponent implements OnInit {
	public pods: Pod[] = [];
	public services: Service[] = [];
	public namespaces: string[] = [];
	public username: string = "";
	private _authService = inject(AuthenticationService);
	private _kubeService = inject(KubeService);

	ngOnInit(): void {
		this.pods = this._kubeService.getAllPods();
		this.services = this._kubeService.getAllServices();
		this.namespaces = this._kubeService.getAllNamespaces();
		this.username = this._authService.getUser().user;
	}

	@ViewChild(ResourceModalComponent)
	modal!: ResourceModalComponent;

	onClick(): void {
		this.modal.open();
	}
}
