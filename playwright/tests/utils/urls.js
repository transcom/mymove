// Base URLs
// change this to 3000 to run against dev instance
const playwrightClientPort = process.env.PLAYWRIGHT_CLIENT_PORT || 4000;
export const milmoveBaseURL = `http://milmovelocal:${playwrightClientPort}`;
export const officeBaseURL = `http://officelocal:${playwrightClientPort}`;
export const adminBaseURL = `http://adminlocal:${playwrightClientPort}`;
