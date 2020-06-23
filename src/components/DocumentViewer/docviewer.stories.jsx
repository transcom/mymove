import React from 'react';
import { storiesOf } from '@storybook/react';

import DocViewerMenu from './Menu';
import DocViewerContent from './Content';

storiesOf('Components|Document Viewer', module)
  .add('menu', () => (
    <div>
      <DocViewerMenu />
    </div>
  ))
  .add('content area', () => (
    <div>
      <DocViewerContent />;
    </div>
  ));
