import React from 'react';
import { text } from '@storybook/addon-knobs';

import SimpleSection from './SimpleSection';

export default {
  title: 'Components/Simple Section',
};

export const Default = () => (
  <SimpleSection header={text('header', 'Simple Section')} border>
    <p>{text('content', 'This is a simple section for simple content.')}</p>
  </SimpleSection>
);

export const WithSubsections = () => (
  <SimpleSection header={text('header', 'Simple Section')} border>
    <SimpleSection header={text('subHeader1', 'Subsection')}>
      <p>Here&apos;s a simple subsection.</p>
      <p>{text('subContent1', 'I have more content here. Maybe I have a lot to say.')}</p>
    </SimpleSection>
    <SimpleSection header={text('subHeader2', 'Another subsection')}>
      <p>{text('subContent2', 'You can have as many as you like.')}</p>
    </SimpleSection>
  </SimpleSection>
);
