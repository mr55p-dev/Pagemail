export interface UserRecord {
	id: string;
	username?: string;
	verified: boolean;
	email?: string;
	created: string;
	updated: string;
	avatar?: string;
	subscribed: boolean;
}

export interface PageRecord {
  id: string;
  url: string;
  user_id: string;
}
