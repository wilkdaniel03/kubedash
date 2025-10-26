import { HttpContextToken } from '@angular/common/http';

export const BEARER_TOKEN_ENABLED = new HttpContextToken<boolean>(() => false);
