/*
import { requiredAddressSchema } from 'utils/validation';
const address1 = {
  city: 'Alexandria',
  state: 'VA',
  postalCode: '12867',
  streetAddress1: '333 Most Fake Blvd',
  streetAddress2: '',
  streetAddress3: '',
  country: 'USA',
};
const defaultValues = {
  closeModal: () => {},
  onSave: () => {},
  isOpen: true,
  serviceItem: dddSitWithAddressUpdate,
};

*/

import React from 'react';

import { ShipmentAddressUpdateReviewRequestModal } from './ShipmentAddressUpdateReviewRequestModal';

export default {
  title: 'Office Components/ShipmentAddressUpdateReviewRequestModal',
  component: ShipmentAddressUpdateReviewRequestModal,
};

export const ReviewModal = {
  render: () => <ShipmentAddressUpdateReviewRequestModal />,
};

/*
export const Approve? = {
  render: () => <ShipmentAddressUpdateReviewRequestModal />,
};

export const Reject? = {
  render: () => <ShipmentAddressUpdateReviewRequestModal />,
};
*/
