import React from 'react';
import moment from 'moment';

import { storiesOf } from '@storybook/react';

import { StatusTimeline } from '../scenes/Landing/StatusTimeline';

const StatusTimelineCodes = {
  Submitted: 'SUBMITTED',
  PpmApproved: 'PPM_APPROVED',
  InProgress: 'IN_PROGRESS',
  PaymentRequested: 'PAYMENT_REQUESTED',
  PaymentReviewed: 'PAYMENT_REVIEWED',
};

const oneStatusBlockComplete = [
  {
    name: 'Submitted',
    code: StatusTimelineCodes.Submitted,
    completed: true,
    dates: [
      moment()
        .subtract(1, 'day')
        .format(),
    ],
  },
  {
    name: 'Approved',
    code: StatusTimelineCodes.PpmApproved,
    completed: false,
  },
  {
    name: 'In progress',
    code: StatusTimelineCodes.InProgress,
    completed: false,
  },
  {
    name: 'Payment requested',
    code: StatusTimelineCodes.PaymentRequested,
    completed: false,
  },
  {
    name: 'Payment reviewed',
    code: StatusTimelineCodes.PaymentReviewed,
    completed: false,
  },
];

const someStatusBlocksCompleted = [
  {
    name: 'Submitted',
    code: StatusTimelineCodes.Submitted,
    completed: true,
    dates: [
      moment()
        .subtract(4, 'days')
        .format(),
    ],
  },
  {
    name: 'Approved',
    code: StatusTimelineCodes.PpmApproved,
    completed: true,
    dates: [
      moment()
        .subtract(2, 'days')
        .format(),
    ],
  },
  {
    name: 'In progress',
    code: StatusTimelineCodes.InProgress,
    completed: true,
    dates: [
      moment()
        .subtract(1, 'day')
        .format(),
    ],
  },
  {
    name: 'Payment requested',
    code: StatusTimelineCodes.PaymentRequested,
    completed: false,
    dates: [moment().format()],
  },
  {
    name: 'Payment reviewed',
    code: StatusTimelineCodes.PaymentReviewed,
    completed: false,
  },
];

const allStatusBlocksCompleted = [
  {
    name: 'Submitted',
    code: StatusTimelineCodes.Submitted,
    completed: true,
    dates: [
      moment()
        .subtract(8, 'days')
        .format(),
    ],
  },
  {
    name: 'Approved',
    code: StatusTimelineCodes.PpmApproved,
    completed: true,
    dates: [
      moment()
        .subtract(6, 'days')
        .format(),
    ],
  },
  {
    name: 'In progress',
    code: StatusTimelineCodes.InProgress,
    completed: true,
    dates: [
      moment()
        .subtract(5, 'days')
        .format(),
    ],
  },
  {
    name: 'Payment requested',
    code: StatusTimelineCodes.PaymentRequested,
    completed: true,
    dates: [
      moment()
        .subtract(4, 'days')
        .format(),
    ],
  },
  {
    name: 'Payment reviewed',
    code: StatusTimelineCodes.PaymentReviewed,
    completed: true,
    dates: [
      moment()
        .subtract(2, 'days')
        .format(),
    ],
  },
];

storiesOf('StatusTimeline/showEstimatedTrue', module)
  .addDecorator(storyFn => <div className="shipment_box_contents">{storyFn()}</div>)
  .add('with one status block complete', () => (
    <StatusTimeline statuses={oneStatusBlockComplete} showEstimated={true} />
  ))
  .add('with more than one but not all status blocks complete', () => (
    <StatusTimeline statuses={someStatusBlocksCompleted} showEstimated={true} />
  ))
  .add('with all status block complete', () => (
    <StatusTimeline statuses={allStatusBlocksCompleted} showEstimated={true} />
  ));

storiesOf('StatusTimeline/showEstimatedFalse', module)
  .addDecorator(storyFn => <div className="shipment_box_contents">{storyFn()}</div>)
  .add('with one status block complete', () => (
    <StatusTimeline statuses={oneStatusBlockComplete} showEstimated={false} />
  ))
  .add('with more than one but not all status blocks complete', () => (
    <StatusTimeline statuses={someStatusBlocksCompleted} showEstimated={false} />
  ))
  .add('with all status block complete', () => (
    <StatusTimeline statuses={allStatusBlocksCompleted} showEstimated={false} />
  ));
