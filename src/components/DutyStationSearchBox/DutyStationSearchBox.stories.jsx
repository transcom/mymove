import React from 'react';

import { DutyStationSearchBox } from 'scenes/ServiceMembers/DutyStationSearchBox';

export default {
  title: 'Components/Duty Station Search Box',
  component: DutyStationSearchBox,
};

const value = {
  address: {
    city: 'Glendale Luke AFB',
    country: 'United States',
    id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
    postal_code: '85309',
    state: 'AZ',
    street_address_1: 'n/a',
  },
  address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
  affiliation: 'AIR_FORCE',
  created_at: '2021-02-11T16:48:04.117Z',
  id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
  name: 'Luke AFB',
  updated_at: '2021-02-11T16:48:04.117Z',
};

export const standard = () => {
  return <DutyStationSearchBox input={{ name: 'test_component' }} title="Test Component" />;
};

export const withValue = () => {
  return (
    <DutyStationSearchBox input={{ name: 'test_component', value }} title="Test Component" displayAddress={false} />
  );
};

export const withValueAndAddress = () => {
  return <DutyStationSearchBox input={{ name: 'test_component', value }} title="Test Component" />;
};

export const withErrorMessage = () => {
  return (
    <DutyStationSearchBox input={{ name: 'test_component' }} title="Test Component" errorMsg="Something went wrong" />
  );
};
