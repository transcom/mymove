import React from 'react';

import SimpleSection from './SimpleSection';

export default {
  title: 'Components/SimpleSection',
};

export const Default = () => (
  <SimpleSection header="Simple Section">
    <p>This is a simple section for simple content.</p>
  </SimpleSection>
);

export const WithSubsections = () => (
  <SimpleSection header="Simple Section">
    <SimpleSection header="Subsection" border={false}>
      <p>Here&apos;s a simple subsection.</p>
      <p>I have more content here. Maybe I have a lot to say.</p>
    </SimpleSection>
    <SimpleSection header="Another subsection" border={false}>
      <p>You can have as many as you like.</p>
    </SimpleSection>
  </SimpleSection>
);
