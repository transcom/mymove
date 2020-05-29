import React from 'react';
import { storiesOf } from '@storybook/react';

import { GridContainer, Grid } from '@trussworks/react-uswds';

storiesOf('Global|Grid', module)
  .add('default container', () => (
    <GridContainer>
      <Grid>Content</Grid>
      <Grid>Content</Grid>
    </GridContainer>
  ))
  .add('column spans', () => (
    <GridContainer>
      <Grid>Content</Grid>
      <Grid>Content</Grid>
    </GridContainer>
  ));
