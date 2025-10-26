import { Component, inject, OnDestroy, OnInit, signal, TemplateRef, ViewChild, ViewContainerRef } from '@angular/core';
import { toObservable } from '@angular/core/rxjs-interop';
import { Select } from '@interfaces';
import { ModalService } from '@services';
import { Subscription } from 'rxjs';

@Component({
	selector: 'resource-modal',
	template: `
		<ng-template #modal>
			<ng-template #deployment><deployment-view/></ng-template>
			<ng-template #service><service-view/></ng-template>
			<layout-modal>
				<h1 class="title">Create resource</h1>
				<core-select [model]="selectModel" [control]="selectState"/>
				<ng-container #content></ng-container>
			</layout-modal>
		</ng-template>
		<ng-container #modalContainer></ng-container>
	`,
	providers: [ModalService]
})
export class ResourceModalComponent implements OnInit, OnDestroy {
	private _modalService = inject(ModalService);

	private status$: Subscription;

	selectModel: Select[] = [
		{ value: 0, content: "Create deployment" },
		{ value: 1, content: "Create service" }
	];

	selectState = signal<string>('');
	selectState$ = toObservable(this.selectState);

	@ViewChild('modalContainer',{read: ViewContainerRef})
	modalContainer!: ViewContainerRef;

	@ViewChild('modal')
	modal!: TemplateRef<any>;

	@ViewChild('content',{read: ViewContainerRef})
	content!: ViewContainerRef;

	@ViewChild('deployment')
	deployment!: TemplateRef<any>;

	@ViewChild('service')
	service!: TemplateRef<any>;

	ngOnInit(): void {
		this.status$ = this._modalService.getStatus().subscribe(status => {
			if(!this.modalContainer) return;
			if(status)
				this.modalContainer.createEmbeddedView(this.modal);
			else
				this.modalContainer.clear();
		});
		this.selectState$.subscribe(res => {
			if(res == this.selectModel[0].content) this.attachDeploymentView();
			else this.attachServiceView();
		})
	}

	private attachDeploymentView(): void {
		this.content.clear();
		this.content.createEmbeddedView(this.deployment);
	}

	private attachServiceView(): void {
		this.content.clear();
		this.content.createEmbeddedView(this.service);
	}

	open(): void {
		this._modalService.open();
	}

	ngOnDestroy(): void {
		this.status$.unsubscribe();
	}
}
