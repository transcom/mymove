import React from 'react';
import { Tag } from '@trussworks/react-uswds';

import WeightDisplay from 'components/Office/WeightDisplay/WeightDisplay';
import { MockProviders } from 'testUtils';
import { permissionTypes } from 'constants/permissions';

export default {
  title: 'Office Components/WeightDisplay',
  component: WeightDisplay,
  argTypes: {
    weightValue: { defaultValue: 10000 },
    onEdit: { defaultValue: null },
    heading: { defaultValue: 'weight allowance' },
    children: { defaultValue: null },
  },
};

const ExternalVendorShipmentMessage = () => (
  <small>
    1 shipment not moved by GHC prime. <a href="">View move details</a>
  </small>
);

const Template = (args) => (
  <MockProviders permissions={[permissionTypes.updateBillableWeight]}>
    <WeightDisplay {...args} />
  </MockProviders>
);

export const WithNoWeight = Template.bind({});
WithNoWeight.args = {
  weightValue: null,
};

export const WithWeight = Template.bind({});

export const WithEditButton = Template.bind({});
WithEditButton.argTypes = {
  onEdit: { defaultValue: () => {}, action: 'clicked' },
};

export const WithWeightAndDetailsTag = Template.bind({});
WithWeightAndDetailsTag.args = {
  children: <Tag>Risk of excess</Tag>,
};

export const WithWeightAndDetailsText = Template.bind({});
WithWeightAndDetailsText.args = {
  children: '110% of estimated weight',
};

export const WithWeightAndExternalVendorNTSRShipment = Template.bind({});
WithWeightAndExternalVendorNTSRShipment.args = {
  children: <ExternalVendorShipmentMessage />,
};

export const WithWeightAndExternalVendorNTSRShipmentAndTag = Template.bind({});
WithWeightAndExternalVendorNTSRShipmentAndTag.args = {
  children: [<Tag>Risk of excess</Tag>, <br />, <ExternalVendorShipmentMessage />],
};
