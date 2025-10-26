import { Component, ElementRef, inject, ViewChild } from '@angular/core';
import { ModalService } from '@services';
import { BehaviorSubject } from 'rxjs';

@Component({
	selector: 'layout-modal',
	template: `
		<div class="wrapper">
			<ng-content></ng-content>
		</div>
		<div class="background" (click)="onClick()"></div>
	`,
	styles: `
		:host {
			display: flex;
			position: fixed;
			width: 100vw;
			height: 100vh;
			top: 0;
			left: 0;
			justify-content: center;
			align-items: center;
		}

		.wrapper {
			width: 40%;
			height: 60%;
			background: #fff;
			border-radius: 25px;
			z-index: 2;
			padding: 15px 20px;
			border: 1px solid #e7e9ec;
		}

		.background {
			position: absolute;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			background: rgba(0,0,0,0.3);
			z-index: 1;
		}
	`
})
export class LayoutModalComponent {
	private modalService = inject(ModalService);
	selected = new BehaviorSubject<number>(0);

	@ViewChild('select')
	select!: ElementRef<HTMLSelectElement>;

	onClick(): void {
		this.modalService.close();
	}
}
