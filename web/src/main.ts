import { TuiRoot,TuiButton } from "@taiga-ui/core";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';
import { BrowserModule } from '@angular/platform-browser';
import { Component, inject, NgModule, OnInit } from '@angular/core';
import { Router as MainRouter, Event, ResolveStart } from '@angular/router';
import { Router } from './router.module';
import { AuthenticationService, HttpService, KubeService } from '@services';
import { TokenInterceptor } from '@interceptors';
import { HTTP_INTERCEPTORS, provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { map, tap } from "rxjs";
import { IPods, IServices } from "@interfaces";

@Component({
	selector: 'app-root',
	template: `<main><router-outlet></router-outlet></main>`,
	styles: `
		main {
			width: 100vw;
			height: 100vh;
		}
	`
})
class App implements OnInit {
	private _router = inject(MainRouter);
	private _http = inject(HttpService);
	private _kubeService = inject(KubeService);

	// DOING HTTP REQUESTS
	ngOnInit(): void {
		this._router.events.subscribe((res: Event) => {
			if(res instanceof ResolveStart && res.url == '/home') {
				this._http.listpods().pipe(
					tap({
						next: (pods: IPods) => {
							for(let pod of pods.pods) {
								this._kubeService.pushPod(pod);
							}
						},
						error: err => {
							console.log("Failed to get pods");
						}
					})
				).subscribe();
				this._http.listservices().pipe(
					tap({
						next: (services: IServices) => {
							for(let service of services.services) {
								this._kubeService.pushService(service);
							}
						},
						error: err => {
							console.log("Failed to get services");
						}
					})
				).subscribe();
				this._http.listnamespaces().pipe(
					tap({
						next: (namespaces: string[]) => {
							for(let ns of namespaces) {
								this._kubeService.pushNamespace(ns);
							}
						},
						error: err => {
							console.log("Failed to get namespaces");
						}
					})
				).subscribe();
			}
		});
	}
}

@NgModule({
	declarations: [App],
	imports: [Router,BrowserModule,BrowserAnimationsModule,TuiRoot,TuiButton],
	providers: [
			AuthenticationService,
			HttpService,
			KubeService,
			provideHttpClient(withInterceptorsFromDi()),
			{ provide: HTTP_INTERCEPTORS, useClass: TokenInterceptor, multi: true}
	],
	bootstrap: [App]
})
class AppModule {}

platformBrowserDynamic().bootstrapModule(AppModule);
