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
	readabilityEnabled: boolean;
}

export interface PageRecord extends BaseRecord {
  id: string;
  url: string;
  user_id: string;
  title?: string;
  description?: string;
  image_url?: string;
  is_readable?: string;
  readability_status?: ReadabilityStatus;
}

export enum ReadabilityStatus {
  UNKNOWN = "UNKNOWN",
  FAILED = "FAILED",
  PROCESSING = "PROCESSING",
  COMPLETE = "COMPLETE",
}
