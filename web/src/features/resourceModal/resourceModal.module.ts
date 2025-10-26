import { NgModule } from '@angular/core';
import { CoreModule } from '@core';
import { LayoutModule } from '@layout';
import { ResourceModalComponent } from './resourceModal.component';
import { ResourceModalDeploymentView } from './resourceModalDeployment.component';
import { ResourceModalServiceView } from './resourceModalService.component';

@NgModule({
	declarations: [
		ResourceModalComponent,
		ResourceModalDeploymentView,
		ResourceModalServiceView
	],
	imports: [CoreModule,LayoutModule],
	exports: [ResourceModalComponent]
})
export class ResourceModalModule {}
