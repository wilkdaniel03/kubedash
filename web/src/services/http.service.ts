import { inject, Injectable } from '@angular/core';
import { HttpClient, HttpContext, } from '@angular/common/http';
import { IUser, ILoginResult, IPods, IServices } from '@interfaces';
import { Observable } from 'rxjs';
import { BEARER_TOKEN_ENABLED } from '@interceptors';

const URL = "http://localhost:5000/"; 

@Injectable()
export class HttpService {
	private _http = inject(HttpClient);

	public login(data: IUser): Observable<ILoginResult> {
		const url = URL + "api/v1/login";
		return this._http.post<ILoginResult>(url,data);
	}

	public listnamespaces(): Observable<string[]> {
		const url = URL + "api/v1/listnamespaces";
		return this._http.get<string[]>(url,{
			context: new HttpContext().set(BEARER_TOKEN_ENABLED,true)
		});
	}

	public listpods(): Observable<IPods> {
		const url = URL + "api/v2/listpods";
		return this._http.get<IPods>(url,{
			context: new HttpContext().set(BEARER_TOKEN_ENABLED,true)
		});
	}

	public listservices(): Observable<IServices> {
		const url = URL + "api/v1/listservices";
		return this._http.get<IServices>(url,{
			context: new HttpContext().set(BEARER_TOKEN_ENABLED,true)
		});
	}

	public createdeployment(namespace: string, name: string, image: string, replicas: number): Observable<any> {
		const url = URL + "api/v1/createdeployment";
		return this._http.post<any>(url,{
			namespace: namespace,
			name: name,
			image: image,
			replicas: replicas
		},{
			context: new HttpContext().set(BEARER_TOKEN_ENABLED,true)
		});
	}
}
