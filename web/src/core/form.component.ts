import { AfterViewInit, Component, ElementRef, EventEmitter, Input, Output, ViewChild, ViewContainerRef, WritableSignal, signal } from '@angular/core';
import { CoreInputComponent, CoreSelectComponent } from '@core';
import { FormField, FORM_FIELD_TYPE, InputType, Select } from '@interfaces';

@Component({
	selector: 'core-form',
	template: `
		<form (submit)="onSubmit($event)">
			<ng-container #content></ng-container>
			<core-btn>{{ submitContent }}</core-btn>
		</form>
	`,
	styles: `
		form {
			height: 100%;
			display: flex;
			flex-direction: column;
			justify-content: space-evenly;
		}
	`
})
export class CoreFormComponent implements AfterViewInit {
	control: Map<string,WritableSignal<string>> = new Map();

	@Input()
	model: FormField[] = [];

	@Input()
	submitContent: string = '';

	@ViewChild('content',{read: ViewContainerRef})
	content!: ViewContainerRef;

	@ViewChild('form')
	form!: ElementRef<HTMLFormElement>;

	@Output()
	submit = new EventEmitter<Record<string,string>>();

	ngAfterViewInit(): void {
		for(let field of this.model) {
			if(field.type == FORM_FIELD_TYPE.INPUT_TEXT || field.type == FORM_FIELD_TYPE.INPUT_PASSWORD || field.type == FORM_FIELD_TYPE.INPUT_NUMBER)
				this.createInput(field);
			else 
				this.createSelect(field);
		}
	}

	private createInput(field: FormField): void {
		let type: InputType;
		if(field.type == FORM_FIELD_TYPE.INPUT_TEXT)
			type = 'text';
		else if(field.type == FORM_FIELD_TYPE.INPUT_NUMBER)
			type = 'number';
		else
			type = 'password';
		const comp = this.content.createComponent(CoreInputComponent);
		comp.instance.type = type;
		comp.instance.placeholder = field.placeholder;
		this.control.set(field.name,signal<string>(''));
		const sig = this.control.get(field.name);
		if(sig) comp.instance.control = sig;
	}

	private createSelect(field: FormField): void {
		if(!field.items) return;
		const obj: Select[] = [];
		for(let i = 0; i < field.items.length; i++) {
			obj.push({value:i,content:field.items[i]});
		}
		const comp = this.content.createComponent(CoreSelectComponent);
		comp.instance.model = obj;
		this.control.set(field.name,signal<string>(''));
		const sig = this.control.get(field.name);
		if(sig) comp.instance.control = sig;
	}

	public onSubmit(event: SubmitEvent): void {
		event.preventDefault();
		const obj: Record<string,string> = {};
		for(let field of this.model) {
			const sig = this.control.get(field.name);
			if(sig) obj[field.name] = sig();
		}
		this.submit.emit(obj);
	}
}
