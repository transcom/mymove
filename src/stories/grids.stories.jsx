import React from 'react';
import { storiesOf } from '@storybook/react';

import { GridContainer, Grid } from '@trussworks/react-uswds';

const exampleStyles = {
  border: '1px solid',
  padding: '1rem',
  backgroundColor: '#e1e7f1',
};

const testContent = <div style={exampleStyles}>Content</div>;

storiesOf('Global|Grid', module)
  .add('default container', () => (
    <GridContainer>
      <Grid row>
        <Grid tablet={{ col: true }}>{testContent}</Grid>
        <Grid tablet={{ col: true }}>{testContent}</Grid>
        <Grid tablet={{ col: true }}>{testContent}</Grid>
      </Grid>
    </GridContainer>
  ))
  .add('column spans', () => (
    <GridContainer>
      <Grid row>
        <Grid col={1}>{testContent}</Grid>
        <Grid col={11}>{testContent}</Grid>
      </Grid>
      <Grid row>
        <Grid col={2}>{testContent}</Grid>
        <Grid col={10}>{testContent}</Grid>
      </Grid>
      <Grid row>
        <Grid col={3}>{testContent}</Grid>
        <Grid col={9}>{testContent}</Grid>
      </Grid>
      <Grid row>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={8}>{testContent}</Grid>
      </Grid>
      <Grid row>
        <Grid col={5}>{testContent}</Grid>
        <Grid col={7}>{testContent}</Grid>
      </Grid>
      <Grid row>
        <Grid col={6}>{testContent}</Grid>
        <Grid col={6}>{testContent}</Grid>
      </Grid>
      <Grid row>
        <Grid col={1}>{testContent}</Grid>
        <Grid col={1}>{testContent}</Grid>
        <Grid col={1}>{testContent}</Grid>
        <Grid col={1}>{testContent}</Grid>
        <Grid col={1}>{testContent}</Grid>
        <Grid col={1}>{testContent}</Grid>
        <Grid col={1}>{testContent}</Grid>
        <Grid col={1}>{testContent}</Grid>
        <Grid col={1}>{testContent}</Grid>
        <Grid col={1}>{testContent}</Grid>
        <Grid col={1}>{testContent}</Grid>
        <Grid col={1}>{testContent}</Grid>
      </Grid>
      <Grid row>
        <Grid col={2}>{testContent}</Grid>
        <Grid col={2}>{testContent}</Grid>
        <Grid col={2}>{testContent}</Grid>
        <Grid col={2}>{testContent}</Grid>
        <Grid col={2}>{testContent}</Grid>
        <Grid col={2}>{testContent}</Grid>
      </Grid>
      <Grid row>
        <Grid col={3}>{testContent}</Grid>
        <Grid col={3}>{testContent}</Grid>
        <Grid col={3}>{testContent}</Grid>
        <Grid col={3}>{testContent}</Grid>
      </Grid>
      <Grid row>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
      </Grid>
    </GridContainer>
  ))
  .add('auto layout columns', () => (
    <GridContainer>
      <Grid row>
        <Grid col="auto">{testContent}</Grid>
        <Grid col="fill">{testContent}</Grid>
        <Grid col="fill">{testContent}</Grid>
        <Grid col="auto">{testContent}</Grid>
      </Grid>
    </GridContainer>
  ))
  .add('responsive', () => (
    <div>
      <h2>Same at all breakpoints</h2>
      <GridContainer>
        <Grid row>
          <Grid col={1}>{testContent}</Grid>
          <Grid col={2}>{testContent}</Grid>
          <Grid col={3}>{testContent}</Grid>
          <Grid col={4}>{testContent}</Grid>
          <Grid col={2}>{testContent}</Grid>
        </Grid>
        <Grid row>
          <Grid col={8}>{testContent}</Grid>
          <Grid col={2}>{testContent}</Grid>
          <Grid col={2}>{testContent}</Grid>
        </Grid>
      </GridContainer>

      <h2>Stacked columns at narrow widths</h2>
      <GridContainer>
        <Grid row>
          <Grid tablet={{ col: true }}>{testContent}</Grid>
          <Grid tablet={{ col: true }}>{testContent}</Grid>
          <Grid tablet={{ col: true }}>{testContent}</Grid>
        </Grid>
        <Grid row>
          <Grid tablet={{ col: 4 }}>{testContent}</Grid>
          <Grid tablet={{ col: 8 }}>{testContent}</Grid>
        </Grid>
      </GridContainer>

      <h2>Mix and match</h2>
      <GridContainer>
        <Grid row>
          <Grid tablet={{ col: 8 }}>{testContent}</Grid>
          <Grid col={6} tablet={{ col: 4 }}>
            {testContent}
          </Grid>
        </Grid>
        <Grid row>
          <Grid col={6} tablet={{ col: 4 }}>
            {testContent}
          </Grid>
          <Grid col={6} tablet={{ col: 4 }}>
            {testContent}
          </Grid>
          <Grid col={6} tablet={{ col: 4 }}>
            {testContent}
          </Grid>
        </Grid>
        <Grid row>
          <Grid col={6}>{testContent}</Grid>
          <Grid col={6}>{testContent}</Grid>
        </Grid>
      </GridContainer>
    </div>
  ))
  .add('offset columns', () => (
    <GridContainer>
      <Grid row>
        <Grid col>{testContent}</Grid>
      </Grid>
      <Grid row>
        <Grid col={1}>{testContent}</Grid>
      </Grid>
      <Grid row>
        <Grid col={1} offset={1}>
          {testContent}
        </Grid>
      </Grid>
      <Grid row>
        <Grid col={1} offset={2}>
          {testContent}
        </Grid>
      </Grid>
      <Grid row>
        <Grid col={1} offset={3}>
          {testContent}
        </Grid>
      </Grid>
      <Grid row>
        <Grid col={1} offset={4}>
          {testContent}
        </Grid>
      </Grid>
      <Grid row>
        <Grid col={1} offset={5}>
          {testContent}
        </Grid>
      </Grid>
      <Grid row>
        <Grid col={1} offset={6}>
          {testContent}
        </Grid>
      </Grid>
      <Grid row>
        <Grid col={1} offset={7}>
          {testContent}
        </Grid>
      </Grid>
      <Grid row>
        <Grid col={1} offset={8}>
          {testContent}
        </Grid>
      </Grid>
      <Grid row>
        <Grid col={1} offset={9}>
          {testContent}
        </Grid>
      </Grid>
      <Grid row>
        <Grid col={1} offset={10}>
          {testContent}
        </Grid>
      </Grid>
      <Grid row>
        <Grid col={1} offset={11}>
          {testContent}
        </Grid>
      </Grid>
      <Grid row>
        <Grid col={8} offset={4}>
          {testContent}
        </Grid>
      </Grid>
    </GridContainer>
  ))
  .add('column wrapping', () => (
    <GridContainer>
      <Grid row>
        <Grid col={8}>{testContent}</Grid>
        <Grid col={3}>{testContent}</Grid>
        <Grid col={5}>{testContent}</Grid>
      </Grid>
    </GridContainer>
  ));
