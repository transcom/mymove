# How to display dates and times

Timezones are hard to do correctly, especially factoring in daylight savings time and users in multiple locations using our app.

## Displaying dates in forms

Other government projects use the date that the person filling the form would have used (their date in their local timezone). This is what we'll continue doing. The official format in the office and TSP app so far has been DD-Mon-YYYY (29-Mar-2018, as an example). In the service member interfaces, we use the date format MM/DD/YYYY (03/29/2018). This should be what we default to unless otherwise stated. We don't currently store a user's local timezone, so unless we make changes to do that, we should send the timezone for the user back to the server from the client. Note that this will get more complicated when we add OCONUS moves, but in the meantime it should suffice.
