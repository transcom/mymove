export const longPageLoadTimeout = 10000;
export const fileUploadTimeout = 10000;

// Base URLs
let cypressClientPort = Cypress.env('testClientPort') || 4000; //change this to 3000 to run cypress against dev instance
export const milmoveBaseURL = `http://milmovelocal:${cypressClientPort}`;
export const officeBaseURL = `http://officelocal:${cypressClientPort}`;

// App Types
export const milmoveAppName = 'milmove';
export const officeAppName = 'office';

// User Types
export const milmoveUserType = 'milmove';
export const PPMOfficeUserType = 'PPM office';
export const TOOOfficeUserType = 'TOO office';
export const TIOOfficeUserType = 'TIO office';
export const dpsUserType = 'dps';

// User Types to Base URLs
/*  security/detect-object-injection */
export const userTypeToBaseURL = {};
userTypeToBaseURL[milmoveUserType] = milmoveBaseURL;
userTypeToBaseURL[PPMOfficeUserType] = officeBaseURL;
userTypeToBaseURL[TOOOfficeUserType] = officeBaseURL;
userTypeToBaseURL[TIOOfficeUserType] = officeBaseURL;
userTypeToBaseURL[dpsUserType] = milmoveBaseURL;
/* eslint-enable security/detect-object-injection */
