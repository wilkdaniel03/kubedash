import { AfterViewInit, Component, ElementRef, inject, Input, OnInit, Renderer2, signal, ViewChild, WritableSignal } from '@angular/core';
import { Select } from '@interfaces';

@Component({
	selector: 'core-select',
	template: `
		<select #select class="form-select">
			@for(opt of model; track opt) {
				<option [value]="opt.value">{{ opt.content }}</option>
			}
		</select>
	`
})
export class CoreSelectComponent implements OnInit, AfterViewInit {
	private _renderer = inject(Renderer2);
	
	@Input()
	model: Select[] = [];

	@Input()
	control: WritableSignal<string> = signal('');

	@ViewChild('select')
	select!: ElementRef<HTMLSelectElement>;

	ngOnInit(): void {
		this.control.set(this.model[0].content);
	}

	ngAfterViewInit(): void {
		this._renderer.listen(this.select.nativeElement,'change',(event: any) => {
			this.control.set(this.model[event.target.value].content);
		});
	}
}
