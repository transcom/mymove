/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { withKnobs, text, array } from '@storybook/addon-knobs';

import Helper from '.';

const title = 'Next step: Add your orders';
const separator = '\n';

export const Basic = () => (
  <div className="grid-container">
    <h3>Bulleted list</h3>
    <Helper
      title={text('Title', title)}
      helpList={array(
        'Help List',
        ['If you have a hard copy, you can take photos of each page', 'If you have a PDF, you can upload that'],
        separator,
      )}
    />
    <br />
    <h3>Plain text</h3>
    <Helper
      title={text('Title', title)}
      description={text(
        'Description',
        'If you have a hard copy, you can take photos of each page. If you have a PDF, you can upload that',
      )}
    />
  </div>
);

export default {
  title: 'Customer Components | Helper',
  decorators: [withKnobs],
};
