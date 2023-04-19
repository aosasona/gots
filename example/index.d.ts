/*
* This file is auto-generated and modified by Gots (https://github.com/aosasona/gots). 
* DO NOT MODIFY THE CONTENT OF THIS FILE
*/

export type Profession = string;

export interface Person {
  first_name: string;
  last_name: string;
  dob: string;
  job?: string;
  CreatedAt: string;
  is_active: boolean;
}

export interface Collection {
  name: string;
  whitelisted_users: {
    first_name: string;
    last_name: string;
    dob: string;
    job?: string;
    CreatedAt: string;
    is_active: boolean;
  }[];
  Lead: {
    first_name: string;
    last_name: string;
    dob: string;
    job?: string;
    CreatedAt: string;
    is_active: boolean;
  };
  collection_tags?: string[];
  admin_id?: number;
}
