import React from 'react';
import { action } from '@storybook/addon-actions';

import { DocsButton, EditButton } from '../components/form';

export default {
  title: 'Components/Icon Button',
};

export const docs = () => (
  <div style={{ padding: '20px', display: 'flex', flexWrap: 'wrap' }}>
    <DocsButton label="My Documents" onClick={action('Docs button clicked')} />
    <DocsButton small label="My Documents" onClick={action('Docs button clicked')} />
    <DocsButton secondary label="My Documents" onClick={action('Docs button clicked')} />
    <DocsButton secondary small label="My Documents" onClick={action('Docs button clicked')} />
    <DocsButton unstyled label="My Documents" onClick={action('Docs button clicked')} />
    <DocsButton disabled label="My Documents" onClick={action('Docs button clicked')} />
  </div>
);
export const edit = () => (
  <div style={{ padding: '20px', display: 'flex', flexWrap: 'wrap' }}>
    <EditButton onClick={action('Edit button clicked')} />
    <EditButton small onClick={action('Edit button clicked')} />
    <EditButton secondary onClick={action('Edit button clicked')} />
    <EditButton secondary small onClick={action('Edit button clicked')} />
    <EditButton unstyled onClick={action('Edit button clicked')} />
    <EditButton disabled onClick={action('Edit button clicked')} />
  </div>
);
