import React from 'react';

import ReviewDetailsCard from './ReviewDetailsCard';
import NeedsReview from './NeedsReview';
import RejectRequest from './RejectRequest';
import AuthorizePayment from './AuthorizePayment';

export default {
  title: 'Office Components/ReviewServiceItems/ReviewDetails',
  component: ReviewDetailsCard,
};

export const ReviewDetailsWithNoValues = () => <ReviewDetailsCard />;

export const ReviewDetailsWithError = () => <ReviewDetailsCard completeReviewError={{ detail: 'THIS IS AN ERROR.' }} />;

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

export const ReviewDetailsWithAuthorizePayment = () => (
  <ReviewDetailsCard acceptedAmount={1234} rejectedAmount={1234} requestedAmount={1234}>
    <AuthorizePayment amount={1234} />
  </ReviewDetailsCard>
);
