import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

@Injectable()
export class ModalService {
	private isOpen = new BehaviorSubject<boolean>(false);

	open(): void {
		this.isOpen.next(true);
	}

	close(): void {
		this.isOpen.next(false);
	}

	getStatus(): Observable<boolean> {
		return this.isOpen.asObservable();
	}
}
