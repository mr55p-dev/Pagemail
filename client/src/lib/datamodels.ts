interface BaseRecord {
  created: string;
  updated: string;

}
export interface UserRecord extends BaseRecord {
	id: string;
	username?: string;
	verified: boolean;
	email?: string;
	created: string;
	updated: string;
	avatar?: string;
	subscribed: boolean;
}

export interface PageRecord extends BaseRecord {
  id: string;
  url: string;
  user_id: string;
}
