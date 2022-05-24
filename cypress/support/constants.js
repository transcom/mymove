export const longPageLoadTimeout = 10000;
export const fileUploadTimeout = 10000;

// Base URLs
let cypressClientPort = Cypress.env('testClientPort') || 4000; //change this to 3000 to run cypress against dev instance
export const milmoveBaseURL = `http://milmovelocal:${cypressClientPort}`;
export const officeBaseURL = `http://officelocal:${cypressClientPort}`;
export const adminBaseURL = `http://adminlocal:${cypressClientPort}`;

// App Types
export const milmoveAppName = 'milmove';
export const officeAppName = 'office';

// User Types
export const milmoveUserType = 'milmove';
export const PPMOfficeUserType = 'PPM office';
export const TOOOfficeUserType = 'TOO office';
export const TIOOfficeUserType = 'TIO office';
export const QAECSROfficeUserType = 'QAECSR office';
export const ServicesCounselorOfficeUserType = 'Services Counselor office';
export const PrimeSimulatorUserType = 'Prime Simulator';
export const dpsUserType = 'dps';

// User Types to Base URLs
/* eslint-disable security/detect-object-injection */
export const userTypeToBaseURL = {};
userTypeToBaseURL[milmoveUserType] = milmoveBaseURL;
userTypeToBaseURL[PPMOfficeUserType] = officeBaseURL;
userTypeToBaseURL[TOOOfficeUserType] = officeBaseURL;
userTypeToBaseURL[TIOOfficeUserType] = officeBaseURL;
userTypeToBaseURL[QAECSROfficeUserType] = officeBaseURL;
userTypeToBaseURL[ServicesCounselorOfficeUserType] = officeBaseURL;
userTypeToBaseURL[dpsUserType] = milmoveBaseURL;
/* eslint-enable security/detect-object-injection */
