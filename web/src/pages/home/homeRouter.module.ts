import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomePageComponent } from './index';
import { CoreModule } from '@core';
import { LayoutModule } from '@layout';
import { ResourceModalModule } from '@features';

const routes: Routes = [
	{ path: "home", component: HomePageComponent },
	{ path: "**", redirectTo: "/home" }
];

@NgModule({
	declarations: [HomePageComponent],
	imports: [CoreModule,LayoutModule,ResourceModalModule,RouterModule.forChild(routes)],
	exports: [RouterModule]
})
export class HomeRouter {}
