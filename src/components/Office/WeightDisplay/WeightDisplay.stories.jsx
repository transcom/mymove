import React from 'react';

import WeightDisplay from 'components/Office/WeightDisplay/WeightDisplay';

export default {
  title: 'Office Components/WeightDisplay',
  component: WeightDisplay,
  argTypes: {
    value: { defaultValue: 10000 },
    onEdit: { action: 'clicked' },
    showEditBtn: { defaultValue: false },
    heading: { defaultValue: 'weight allowance' },
  },
};

const Template = (args) => <WeightDisplay {...args} />;

export const WithNoDetails = Template.bind({});
WithNoDetails.args = {
  value: null,
};

export const WithDetails = Template.bind({});

export const WithEditButton = Template.bind({});
WithEditButton.args = {
  showEditBtn: true,
};
