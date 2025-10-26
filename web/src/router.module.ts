import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthPageComponent } from './pages/auth/login.component';
import { HomePageComponent } from './pages/home';
import { CoreModule } from './core';
import { ReactiveFormsModule } from '@angular/forms';
import { LoginFormCfg } from './pages/auth/loginFormConfig.directive';
import { AuthGuard } from './auth.guard';

const loadNew = () => import('./pages/home/homeRouter.module').then(m => m.HomeRouter);

const routes: Routes = [
	{ path: "auth", children: [ { path: "login", component: AuthPageComponent }, { path: "**", redirectTo: "/auth/login"} ] },
	{ path: "home", component: HomePageComponent, canActivate: [AuthGuard] },
	{ path: "**", redirectTo: "/auth" }
]

@NgModule({
	declarations: [AuthPageComponent,LoginFormCfg],
	imports: [CoreModule,ReactiveFormsModule,RouterModule.forRoot(routes)],
	exports: [RouterModule]
})
export class Router {}
