export const longPageLoadTimeout = 10000;
export const fileUploadTimeout = 10000;

// Base URLs
const cypressClientPort = 4000; //change this to 3000 to run cypress against dev instance
export const milmoveBaseURL = `http://milmovelocal:${cypressClientPort}`;
export const officeBaseURL = `http://officelocal:${cypressClientPort}`;
export const tspBaseURL = `http://tsplocal:${cypressClientPort}`;

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
