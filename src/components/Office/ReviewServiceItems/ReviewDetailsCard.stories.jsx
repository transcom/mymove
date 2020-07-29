import React from 'react';

import ReviewDetailsCard from './ReviewDetailsCard';
import NeedsReview from './NeedsReview';
import RejectRequest from './RejectRequest';

export default {
  title: 'TOO/TIO Components|ReviewServiceItems/ReviewDetails',
  component: ReviewDetailsCard,
};

export const ReviewDetailsWithNoValues = () => <ReviewDetailsCard />;

export const ReviewDetailsWithNeedsReview = () => (
  <ReviewDetailsCard acceptedAmount={1234} rejectedAmount={1234} requestedAmount={1234}>
    <NeedsReview numberOfItems={1} />
  </ReviewDetailsCard>
);

export const ReviewDetailsWithRejectRequest = () => (
  <ReviewDetailsCard acceptedAmount={1234} rejectedAmount={1234} requestedAmount={1234}>
    <RejectRequest />
  </ReviewDetailsCard>
);
