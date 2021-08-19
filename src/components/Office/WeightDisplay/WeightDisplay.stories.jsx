import React from 'react';
import { Tag } from '@trussworks/react-uswds';

import WeightDisplay from 'components/Office/WeightDisplay/WeightDisplay';

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

const Template = (args) => <WeightDisplay {...args} />;

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
