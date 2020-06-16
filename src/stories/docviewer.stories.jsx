import React from 'react';
import { storiesOf } from '@storybook/react';

import DocViewerMenu from '../components/DocumentViewer/menu';
import DocViewerContent from '../components/DocumentViewer/content';

storiesOf('Components|Document Viewer', module).add('menu', () => (
  <div className="display-flex">
    <DocViewerMenu />
    <DocViewerContent />
  </div>
));
