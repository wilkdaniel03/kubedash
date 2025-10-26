import { AfterViewInit, Component, ElementRef, Input, signal, ViewChild, WritableSignal } from "@angular/core";
import type { InputType } from '@interfaces';

@Component({
	selector: 'core-input',
	template: `
		<input #input (keyup)="onChange($event)" class="form-control" placeholder="{{placeholder}}" [type]="type" [value]="control()"/>
	`,
	styles: `
		input {
			width: 100%;
			padding: 10px 0px 10px 12px;
		}
	`
})
export class CoreInputComponent implements AfterViewInit {
	@Input()
	placeholder?: string;

	@Input()
	type: InputType = "text";

	@Input()
	rounded: boolean = false;

	@Input()
	control: WritableSignal<string> = signal('');

	@ViewChild('input')
	input!: ElementRef<HTMLInputElement>;

	onChange(e: KeyboardEvent) {
		const value = (e.target as HTMLInputElement).value;
		this.control.set(value);
	}

	ngAfterViewInit(): void {
		if(this.rounded) {
			this.input.nativeElement.style.borderRadius = "25px";
			this.input.nativeElement.style.padding = "10px 20px";
			this.input.nativeElement.style.borderRadius = "25px";
			this.input.nativeElement.style.boxShadow = "0 1px 12px 0 #e7e9ea";
		}
	}
}
