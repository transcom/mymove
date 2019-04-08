export const longPageLoadTimeout = 10000;
export const fileUploadTimeout = 10000;

// Base URLs
export const milmoveBaseURL = 'http://milmovelocal:4000';
export const officeBaseURL = 'http://officelocal:4000';
export const tspBaseURL = 'http://tsplocal:4000';

// App Types
export const milmoveAppName = 'milmove';
export const officeAppName = 'office';
export const tspAppName = 'tsp';

// User Types
export const milmoveUserType = 'milmove';
export const officeUserType = 'office';
export const tspUserType = 'tsp';
export const dpsUserType = 'dps';

// User Types to Base URLs
/* eslint-disable security/detect-object-injection */
export const userTypeToBaseURL = {};
userTypeToBaseURL[milmoveUserType] = milmoveBaseURL;
userTypeToBaseURL[officeUserType] = officeBaseURL;
userTypeToBaseURL[tspUserType] = tspBaseURL;
userTypeToBaseURL[dpsUserType] = milmoveBaseURL;
/* eslint-enable security/detect-object-injection */
