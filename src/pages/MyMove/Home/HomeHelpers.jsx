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
  <Helper title="Time for step 3: Set up your shipments" className={styles['helper-paragraph-only']}>
    <p>Share where and when you&apos;re moving, and how you want your things to be shipped.</p>
    <p>
      Important Notice: USTRANSCOM has contracted a single move manager to manage the hundreds of commercial moving
      companies that pack, ship, and deliver personal property worldwide. They will manage household goods,
      storage-in-transit warehouse services, and unaccompanied baggage shipments. This move manager will be your primary
      contact for scheduling and conducting your move. The DOD will oversee their work and your local transportation
      office will remain your primary DOD contact to ensure quality performance.
    </p>
    <p>
      (Note: This move manager does not manage personally procured moves or replace the existing programs for
      non-temporary storage or the movement of privately owned vehicles.)
    </p>
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
        <strong>We have assigned you a move code above.</strong> Write it down and use this code when talking to any
        representative about your move. You will also receive this code via a confirmation email.
      </p>
    </div>
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
      <p>
        You can start packing, but do not move any of your personal property until you hear that your move is approved
      </p>
    </div>
    <div>
      <p>
        <strong>For HHGs and other shipments using movers</strong>
      </p>
      <div className={styles['top-gap']}>
        <p>Your Customer Care Representative will contact you to:</p>
        <ul>
          <li>Estimate the weight of your personal property</li>
          <li>Finalize dates to pack and pick up your personal property</li>
        </ul>
      </div>
    </div>
    <div>
      <p>
        <strong>
          We recommend visiting&nbsp;
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
        When you are done moving your things, select <strong>Upload PPM Documents</strong> to document your PPM,
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
  <Helper title="Next step: Contact your movers (if you have them)" className={styles['helper-paragraph-only']}>
    <p>
      If your destination changed or your move was canceled, contact your movers ASAP to let them know. They&apos;ll
      work with you to coordinate any changes.
    </p>
  </Helper>
);

export const HelperPPMCloseoutSubmitted = () => (
  <Helper title="Someone will review all of your PPM documentation" className={styles['helper-paragraph-only']}>
    <p>
      If your documentation is clear and valid, you’ll be able to download a payment packet. You can submit that packet
      to your Finance office to finalize your acutal incentive amount and request payment.
    </p>
    <p>
      If any documentation is unclear or inaccurate, the counselor will reach out to you to request clarification or
      documentation updates.
    </p>
  </Helper>
);
