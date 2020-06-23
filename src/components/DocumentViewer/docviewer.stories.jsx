import React from 'react';
import { storiesOf } from '@storybook/react';

import DocViewerMenu from './Menu';
import DocViewerContent from './Content';

storiesOf('Components|Document Viewer|Menu', module).add('menu', () => (
  <div>
    <DocViewerMenu />
  </div>
));

storiesOf('Components|Document Viewer|Content', module).add('content area', () => (
  <div>
    <DocViewerContent />;
  </div>
));
