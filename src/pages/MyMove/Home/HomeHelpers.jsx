import React from 'react';
import { Link } from '@trussworks/react-uswds';

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

export const HelperSubmittedPPM = () => (
  <Helper title="For your do-it-yourself shipments (PPMs)">
    <ul>
      {/* TBD, add associated link to the text 'Visit the MilMove PPM page' */}
      <li>Visit the MilMove PPM page to learn more about DITY shipments and to manage yours.</li>
      <li>Once you have moved, you’ll request payment using MilMove.</li>
    </ul>
  </Helper>
);

export const HelperSubmittedNoPPM = () => (
  <Helper title="What’s next?">
    <ul>
      <li>
        Create a custom checklist at{' '}
        <Link
          variant="external"
          target="_blank"
          rel="noopener noreferrer"
          href="https://planmymove.militaryonesource.mil/"
        >
          Plan My Move
        </Link>
      </li>
      <li>
        <Link
          variant="external"
          target="_blank"
          rel="noopener noreferrer"
          href="https://installations.militaryonesource.mil/"
        >
          Learn more
        </Link>{' '}
        about your new duty station
      </li>
      <li>If details about your move change, talk to your move counselor or with your movers</li>
      <li>Your movers will help you estimate the weight of your things</li>
      <li>They’ll also finalize moving dates</li>
    </ul>
  </Helper>
);
