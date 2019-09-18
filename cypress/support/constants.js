export const longPageLoadTimeout = 10000;
export const fileUploadTimeout = 10000;

// Base URLs
const cypressClientPort = 4000; //change this to 3000 to run cypress against dev instance
export const milmoveBaseURL = `http://milmovelocal:${cypressClientPort}`;
export const officeBaseURL = `http://officelocal:${cypressClientPort}`;

// App Types
export const milmoveAppName = 'milmove';
export const officeAppName = 'office';

// User Types
export const milmoveUserType = 'milmove';
export const officeUserType = 'office';
export const dpsUserType = 'dps';

// User Types to Base URLs
/* eslint-disable security/detect-object-injection */
export const userTypeToBaseURL = {};
userTypeToBaseURL[milmoveUserType] = milmoveBaseURL;
userTypeToBaseURL[officeUserType] = officeBaseURL;
userTypeToBaseURL[dpsUserType] = milmoveBaseURL;
/* eslint-enable security/detect-object-injection */
