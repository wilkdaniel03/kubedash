import { NgModule } from '@angular/core';
import {
	LayoutNavbarComponent,
	LayoutHeaderComponent,
	LayoutModalComponent
} from './index';
import { CoreModule } from '@core';

@NgModule({
	declarations: [
		LayoutNavbarComponent,
		LayoutHeaderComponent,
		LayoutModalComponent
	],
	imports: [CoreModule],
	exports: [
		LayoutNavbarComponent,
		LayoutHeaderComponent,
		LayoutModalComponent
	]
})
export class LayoutModule {}
