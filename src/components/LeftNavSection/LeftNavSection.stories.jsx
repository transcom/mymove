import React from 'react';
import { action } from '@storybook/addon-actions';

import LeftNavSection from './LeftNavSection';

export default {
  title: 'Components/Left Nav Section',
  component: LeftNavSection,
};

export const Basic = () => (
  <LeftNavSection sectionName="testSection" onClickHandler={action('clicked the test section')}>
    Test Section
  </LeftNavSection>
);
