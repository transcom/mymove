import React from 'react';

import Helper from 'components/Customer/Home/Helper';

export const HelperNeedsOrders = () => (
  <Helper title="Next step: Add your orders">
    <ul>
      <li>If you have a hard copy, you can take photos of each page</li>
      <li>If you have a PDF, you can upload that</li>
    </ul>
  </Helper>
);

export const HelperNeedsShipment = () => (
  <Helper title="Gather this info, then plan your shipments">
    <ul>
      <li>Preferred moving details</li>
      <li>Destination address (your new place, your duty station ZIP, or somewhere else)</li>
      <li>Names and contact info for anyone you authorize to act on your behalf</li>
    </ul>
  </Helper>
);

export const HelperNeedsSubmitMove = () => (
  <Helper title="Time to submit your move">
    <ul>
      <li>Double check the info you’ve entered</li>
      <li>Sign the legal agreement</li>
      <li>You’ll hear from a move counselor or your transportation office within a few days</li>
    </ul>
  </Helper>
);

export const HelperSubmittedMove = () => (
  <Helper title="Track your HHG move here">
    <ul>
      <li>Create a custom checklist at Plan My Move</li>
      <li>Learn more about your new duty station</li>
    </ul>
  </Helper>
);
