import React from 'react';

import ToolTip from './ToolTip';

const storybookExport = {
  title: 'Office Components/ToolTip',
  component: ToolTip,
};

export default storybookExport;

const Template = (args) => {
  return <ToolTip {...args} />;
};

export const Basic = Template.bind({});
Basic.args = {
  text: 'Lorem Ipsum',
  style: { margin: '75px 0 0 100px' },
};

export const AlternateIcon = Template.bind({});
AlternateIcon.args = {
  icon: 'info-circle',
  text: 'Lorem Ipsum',
  style: { margin: '75px 0 0 100px' },
};

export const BottomText = Template.bind({});
BottomText.args = {
  position: 'bottom',
  text: 'Lorem Ipsum',
  style: { margin: '0 0 75px 100px' },
};

export const LeftText = Template.bind({});
LeftText.args = {
  position: 'left',
  text: 'Lorem Ipsum',
  style: { margin: '25px 0 25px 225px' },
};

export const RightText = Template.bind({});
RightText.args = {
  position: 'right',
  text: 'Lorem Ipsum',
  style: { margin: '25px 0 25px 0' },
};

export const LargeTextArea = Template.bind({});
LargeTextArea.args = {
  textAreaSize: 'large',
  text: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.',
  style: { margin: '175px 0 0 150px' },
};
