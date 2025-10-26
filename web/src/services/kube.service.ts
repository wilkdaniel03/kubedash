import { Injectable } from '@angular/core';
import { Pod, Service } from '@interfaces';

@Injectable()
export class KubeService {
	private pods: Pod[] = [];
	private services: Service[] = [];
	private namespaces: string[] = [];

	public getAllPods(): Pod[] {
		return this.pods;
	}

	public pushPod(data: Pod): void {
		this.pods.push(data);
	}

	public getAllServices(): Service[] {
		return this.services;
	}

	public pushService(data: Service): void {
		this.services.push(data);
	}

	public getAllNamespaces(): string[] {
		return this.namespaces;
	}

	public pushNamespace(ns: string): void {
		this.namespaces.push(ns);
	}
}
