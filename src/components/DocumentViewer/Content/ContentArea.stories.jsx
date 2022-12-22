import React from 'react';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import pdf from '../sample.pdf';
import xls from '../sample.xls';

import DocViewerContent from './Content';

export default {
  title: 'Components/Document Viewer/Content',
  component: DocViewerContent,
  parameters: {
    happo: false,
  },
  decorators: [
    (Story) => (
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <Story />
          </Grid>
        </Grid>
      </GridContainer>
    ),
  ],
};

const Template = (args) => <DocViewerContent {...args} />;

export const PDF = Template.bind({});
PDF.args = {
  fileType: 'pdf',
  filePath: pdf,
};

export const XLS = Template.bind({});
XLS.args = {
  fileType: 'xls',
  filePath: xls,
};
