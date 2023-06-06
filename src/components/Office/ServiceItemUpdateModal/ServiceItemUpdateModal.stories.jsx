import React from 'react';

import { ServiceItemUpdateModal } from './ServiceItemUpdateModal';
import EditSitAddressChangeForm from './EditSitAddressChangeForm';
import ReviewSitAddressChange from './ReviewSitAddressChange';
import { dddSitWithAddressUpdate } from './ServiceItemUpdateModalTestParams';

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
export default {
  title: 'Office Components/ServiceItemUpdateModal',
  component: ServiceItemUpdateModal,
};

// Story for base component of the Modal
export const ServiceItemUpdateModalStory = {
  render: () => <ServiceItemUpdateModal title="Base modal" {...defaultValues} />,
};
// Story for Editing service item address
export const EditServiceItemAddress = {
  render: () => (
    <ServiceItemUpdateModal
      initialValues={{ officeRemarks: '', newAddress: address1 }}
      validations={{ newAddress: requiredAddressSchema }}
      title="Edit service item"
      {...defaultValues}
    >
      <EditSitAddressChangeForm initialAddress={address1} />
    </ServiceItemUpdateModal>
  ),
};
// Story for reviewing service item addres update request
export const ReviewSitAddress = {
  render: () => (
    <ServiceItemUpdateModal title="Review request: service item update" {...defaultValues}>
      <ReviewSitAddressChange sitAddressUpdate={defaultValues.serviceItem.sitAddressUpdates[0]} />
    </ServiceItemUpdateModal>
  ),
};
