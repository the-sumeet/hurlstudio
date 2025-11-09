// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}

		// Hurl JSON Report Schema
		export interface HurlReport {
			cookies: Cookie[];
			entries: Entry[];
			filename: string;
			success: boolean;
			time: number;
		}

		export interface Entry {
			asserts: Assert[];
			calls: Call[];
			captures: Capture[];
			curl_cmd: string;
			index: number;
			line: number;
			time: number;
		}

		export interface Assert {
			line: number;
			success: boolean;
			message?: string;
		}

		export interface Call {
			request: Request;
			response: Response;
			timings: Timings;
		}

		export interface Request {
			cookies: Cookie[];
			headers: Header[];
			method: string;
			query_string: QueryParam[];
			url: string;
		}

		export interface Response {
			body: string;
			certificate?: Certificate;
			cookies: Cookie[];
			headers: Header[];
			http_version: string;
			status: number;
		}

		export interface Certificate {
			expire_date: string;
			issuer: string;
			serial_number: string;
			start_date: string;
			subject: string;
		}

		export interface Header {
			name: string;
			value: string;
		}

		export interface Cookie {
			name: string;
			value: string;
			domain?: string;
			path?: string;
			expires?: string;
			max_age?: number;
			http_only?: boolean;
			secure?: boolean;
			same_site?: string;
		}

		export interface QueryParam {
			name: string;
			value: string;
		}

		export interface Capture {
			name: string;
			value: any;
		}

		export interface Timings {
			app_connect: number;
			begin_call: string;
			connect: number;
			end_call: string;
			name_lookup: number;
			pre_transfer: number;
			start_transfer: number;
			total: number;
		}
	}
}

export {};
