import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { ReactiveFormsModule } from '@angular/forms';
import {
	CoreInputComponent,
	CoreButtonComponent,
	CoreFormComponent,
	CoreDotComponent,
	CoreLinkComponent,
	CoreSelectComponent
} from './index';

@NgModule({
	declarations: [
		CoreInputComponent,
		CoreButtonComponent,
		CoreFormComponent,
		CoreDotComponent,
		CoreLinkComponent,
		CoreSelectComponent
	],
	imports: [RouterModule,ReactiveFormsModule],
	exports: [
		CoreInputComponent,
		CoreButtonComponent,
		CoreFormComponent,
		CoreDotComponent,
		CoreLinkComponent,
		CoreSelectComponent
	]
})
export class CoreModule {}
