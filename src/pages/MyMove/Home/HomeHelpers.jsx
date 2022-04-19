import React from 'react';
import { Link } from 'react-router-dom';
import { Link as ExternalLink } from '@trussworks/react-uswds';

import styles from './Home.module.scss';

import { customerRoutes } from 'constants/routes';
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
  <Helper title="Time for step 3: set up your shipments">
    <p>Share where and when you&apos;re moving, and how you want your things to be shipped.</p>
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
  <Helper title="Next step: Your move gets approved" className={styles['helper-submitted-section']}>
    <div>
      <p>
        <strong>A move counselor will contact you.</strong> They will confirm the information you entered here, give
        advice, and answer questions.
      </p>
    </div>
    <div>
      <p>
        <strong>For PPM (do it yourself) shipments</strong>
      </p>
      <ul className={styles['top-gap']}>
        <li>You can start packing, but do not move any of your things until you hear that your move is approved</li>
      </ul>
    </div>
    <div>
      <p>
        <strong>For HHGs and other shipments using movers</strong>
      </p>
      <div className={styles['top-gap']}>
        <p>Your movers will contact you to:</p>
        <ul>
          <li>Estimate the weight of your belongings</li>
          <li>Finalize dates to pack and pick up your things</li>
        </ul>
      </div>
    </div>
    <div>
      <p>
        <strong>
          We recommend visiting &nbsp;
          <ExternalLink
            variant="external"
            target="_blank"
            rel="noopener noreferrer"
            href="https://planmymove.militaryonesource.mil/"
          >
            Plan My Move
          </ExternalLink>{' '}
          to make a customized moving checklist.
        </strong>
      </p>
    </div>
  </Helper>
);

export const HelperApprovedMove = () => (
  <Helper title="Your move is in progress." className={styles['helper-approved-section']}>
    <div>
      <p>Talk to your counselor or to your movers to make any changes to your move.</p>
    </div>
    <div>
      <p>
        <strong>For PPM shipments</strong>
      </p>
      <p className={styles['top-gap']}>
        When you are done moving your things, select <strong>Upload PPM documents</strong> to document your PPM,
        calculate your final incentive, and create a payment request packet. You will upload weight tickets, receipts,
        and other documentation that a counselor will review.
      </p>
    </div>
    <div>
      <p>
        <strong>If you receive new orders while your move is underway</strong>
      </p>
      <ul className={styles['top-gap']}>
        <li>Talk to your counselor</li>
        <li>Talk to your movers</li>
        <li>
          <Link to={customerRoutes.ORDERS_AMEND_PATH}>Upload a copy of your new orders</Link>
        </li>
      </ul>
    </div>
  </Helper>
);

export const HelperAmendedOrders = () => (
  <Helper title="Next step: Contact your movers (if you have them)">
    <p>
      If your destination changed or your move was canceled, contact your movers ASAP to let them know. They&apos;ll
      work with you to coordinate any changes.
    </p>
  </Helper>
);
