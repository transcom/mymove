import React from 'react';
import moment from 'moment';
import { withKnobs, boolean, text, date } from '@storybook/addon-knobs';

import { StatusTimeline } from '../scenes/PpmLanding/StatusTimeline';

const StatusTimelineCodes = {
  Submitted: 'SUBMITTED',
  PpmApproved: 'PPM_APPROVED',
  InProgress: 'IN_PROGRESS',
  PaymentRequested: 'PAYMENT_REQUESTED',
  PaymentReviewed: 'PAYMENT_REVIEWED',
};

export default {
  title: 'scenes|Landing',
  decorators: [withKnobs, (storyFn) => <div className="shipment_box_contents">{storyFn()}</div>],
};

export const StatusTimeLine = () => (
  <StatusTimeline
    showEstimated={boolean('showEstimated', false)}
    statuses={[
      {
        name: text('StatusBlock1.Name', 'Submitted', 'Block 1'),
        code: StatusTimelineCodes.Submitted,
        completed: boolean('StatusBlock1.Completed', true, 'Block 1'),
        dates: [date('StatusBlock1.Date1', moment('2020-02-02').subtract(8, 'days').toDate(), 'Block 1')],
      },
      {
        name: text('StatusBlock2.Name', 'Approved', 'Block 2'),
        code: StatusTimelineCodes.PpmApproved,
        completed: boolean('StatusBlock2.Completed', false, 'Block 2'),
        dates: [date('StatusBlock2.Date1', moment('2020-02-02').subtract(6, 'days').toDate(), 'Block 2')],
      },
      {
        name: text('StatusBlock3.Name', 'In Progress', 'Block 3'),
        code: StatusTimelineCodes.InProgress,
        completed: boolean('StatusBlock3.Completed', false, 'Block 3'),
        dates: [date('StatusBlock3.Date1', moment('2020-02-02').subtract(5, 'days').toDate(), 'Block 3')],
      },
      {
        name: text('StatusBlock4.Name', 'Payment Requested', 'Block 4'),
        code: StatusTimelineCodes.PaymentRequested,
        completed: boolean('StatusBlock4.Completed', false, 'Block 4'),
        dates: [date('StatusBlock4.Date1', moment('2020-02-02').subtract(4, 'days').toDate(), 'Block 4')],
      },
      {
        name: text('StatusBlock5.Name', 'Payment Reviewed', 'Block 5'),
        code: StatusTimelineCodes.PaymentReviewed,
        completed: boolean('StatusBlock5.Completed', false, 'Block 5'),
        dates: [date('StatusBlock5.Date1', moment('2020-02-02').subtract(2, 'days').toDate(), 'Block 5')],
      },
    ]}
  />
);
