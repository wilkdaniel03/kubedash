import { AfterViewInit, Component, inject, OnInit, ViewChild } from '@angular/core';
import { CoreFormComponent } from '@core';
import { FORM_FIELD_TYPE, FormField } from '@interfaces';
import { HttpService, KubeService } from '@services';
import { tap } from 'rxjs';

@Component({
	selector: 'deployment-view',
	template: `<div><core-form [model]="model" submitContent="Create deployment"/></div>`,
	styles: `
		div {
			height: 60%;
		}
	`
})
export class ResourceModalDeploymentView implements OnInit, AfterViewInit {
	private _kubeService = inject(KubeService);
	private _httpService = inject(HttpService);

	model: FormField[] = [
		{ type: FORM_FIELD_TYPE.SELECT, name: "namespace", placeholder: "Namespace", items: [] },
		{ type: FORM_FIELD_TYPE.INPUT_TEXT, name: "name", placeholder: "Name" },
		{ type: FORM_FIELD_TYPE.INPUT_TEXT, name: "image", placeholder: "Image" },
		{ type: FORM_FIELD_TYPE.INPUT_NUMBER, name: "replicas", placeholder: "Replicas" }
	];

	@ViewChild(CoreFormComponent)
	form!: CoreFormComponent;

	ngOnInit(): void {
		this.populateNamespaceField();
	}

	ngAfterViewInit(): void {
		this.form.submit.subscribe(res => {
			this._httpService.createdeployment(res['namespace'],res['name'],res['image'],parseInt(res['replicas']))
				.pipe(tap({
					next: res => { console.log(res); },
					error: res => { console.log(res); }
				})).subscribe();
		});
	}

	private populateNamespaceField(): void {
		const namespaces = this._kubeService.getAllNamespaces()
		if(!this.model[0].items) return;
		for(let ns of namespaces)
			this.model[0].items.push(ns);
	}
}
