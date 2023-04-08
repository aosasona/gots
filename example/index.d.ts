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
  createdAt: string;
  is_active: boolean;
}

export interface Collection {
  name: string;
  people: {
    first_name: string;
    last_name: string;
    dob: string;
    job?: string;
    createdAt: string;
    is_active: boolean;
  }[];
  lead: {
    first_name: string;
    last_name: string;
    dob: string;
    job?: string;
    createdAt: string;
    is_active: boolean;
  };
  tags: string[];
}
