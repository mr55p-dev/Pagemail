export interface UserRecord {
	id: string;
	username?: string;
	verified: boolean;
	email?: string;
	created: string;
	updated: string;
	avatar?: string;
}


export interface PageRecord {
  id: string;
  url: string;
  user_id: string;
}
