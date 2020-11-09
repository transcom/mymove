import React from 'react';
import { Link } from 'react-router-dom';
import { Link as ExternalLink } from '@trussworks/react-uswds';

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
  <Helper title="Next: Talk to a move counselor">
    <p>They’ll contact you soon to let you know what to expect and to answer questions.</p>
    <p>
      <strong>Tell them or your movers if anything changes about your move.</strong>
    </p>
    <p>
      <strong>If you’re using government-funded movers, they’ll contact you soon to:</strong>
    </p>
    <ul>
      <li>estimate the weight of your move</li>
      <li>finalize packing and pickup dates</li>
    </ul>
    <p>
      <strong>For more moving tips:</strong>
    </p>
    <ul>
      <li>
        <ExternalLink
          variant="external"
          target="_blank"
          rel="noopener noreferrer"
          href="https://planmymove.militaryonesource.mil/"
        >
          Create a custom checklist at
        </ExternalLink>{' '}
        Plan My Move
      </li>
      <li>
        <ExternalLink
          variant="external"
          target="_blank"
          rel="noopener noreferrer"
          href="https://installations.militaryonesource.mil/"
        >
          Learn more
        </ExternalLink>{' '}
        about your new duty station
      </li>
    </ul>
  </Helper>
);

export const HelperSubmittedPPM = () => (
  <Helper title="For your do-it-yourself shipments (PPMs)">
    <ul>
      <li>
        <Link to="/ppm">Visit the MilMove PPM page</Link> to learn more about DITY shipments and to manage yours.
      </li>
      <li>Once you have moved, you’ll request payment using MilMove.</li>
    </ul>
  </Helper>
);
