export interface IUser {
	user: string;
	pass: string;
}

export enum AUTH_STATUS {
	PENDING,
	OK,
	FAILED
}

export interface ILoginResult {
	token: string;
}

export interface ApiError {
	error: string;
	param?: string;
}

export interface Pod {
	name: string;
	namespace: string;
	status: string;
}

/**
  * interface used in http request to get pods
*/
export interface IPods {
	pods: Pod[]
}

export interface Service {
	name: string;
	namespace: string;
	ports: { name: string; protocol: 'TCP' | 'UDP'; port: number; target_port: number; }[];
	selector: string;
	type: string;
}

/**
  * interface used in http request to get services
*/
export interface IServices {
	services: Service[];
}

export interface Select {
	value: number;
	content: string;
}

export type InputType = "text" | "number" | "password";

export interface IInputConfig {
	label: string;
	type: InputType;
	placeholder?: string;
}

export enum FORM_FIELD_TYPE {
	INPUT_TEXT,
	INPUT_PASSWORD,
	INPUT_NUMBER,
	SELECT
}

export interface FormField {
	type: FORM_FIELD_TYPE,
	items?: string[],
	name: string,
	placeholder?: string
}
