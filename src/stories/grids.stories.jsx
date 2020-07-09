import React from 'react';
import { GridContainer, Grid } from '@trussworks/react-uswds';

const exampleStyles = {
  border: '1px solid',
  padding: '1rem',
  backgroundColor: '#e1e7f1',
};

const testContent = <div style={exampleStyles}>Content</div>;

export default {
  title: 'Global|Grid',
};

export const DefualtContainer = () => (
  <GridContainer>
    <Grid row>
      <Grid tablet={{ col: true }}>{testContent}</Grid>
      <Grid tablet={{ col: true }}>{testContent}</Grid>
      <Grid tablet={{ col: true }}>{testContent}</Grid>
    </Grid>
  </GridContainer>
);
export const ColumnSpans = () => (
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
);
export const AutoLayoutColumns = () => (
  <GridContainer>
    <Grid row>
      <Grid col="auto">{testContent}</Grid>
      <Grid col="fill">{testContent}</Grid>
      <Grid col="fill">{testContent}</Grid>
      <Grid col="auto">{testContent}</Grid>
    </Grid>
  </GridContainer>
);
export const Responsive = () => (
  <div>
    <h4>Same at all breakpoints</h4>
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

    <h4>Stacked columns at narrow widths</h4>
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

    <h4>Mix and match</h4>
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
);
export const OffsetColumns = () => (
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
);
export const ColumnWrapping = () => (
  <GridContainer>
    <Grid row>
      <Grid col={8}>{testContent}</Grid>
      <Grid col={3}>{testContent}</Grid>
      <Grid col={5}>{testContent}</Grid>
    </Grid>
  </GridContainer>
);
export const Gutters = () => (
  <div>
    <h4>Default gutter</h4>
    <GridContainer>
      <Grid row gap>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
      </Grid>
    </GridContainer>

    <h4>Small gutter</h4>
    <GridContainer>
      <Grid row gap="sm">
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
      </Grid>
    </GridContainer>

    <h4>Medium gutter</h4>
    <GridContainer>
      <Grid row gap="md">
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
      </Grid>
    </GridContainer>

    <h4>Large gutter</h4>
    <GridContainer>
      <Grid row gap="lg">
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
      </Grid>
    </GridContainer>

    <h4>2px gutter</h4>
    <GridContainer>
      <Grid row gap="2px">
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
      </Grid>
    </GridContainer>

    <h4>05 gutter</h4>
    <GridContainer>
      <Grid row gap="05">
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
      </Grid>
    </GridContainer>

    <h4>1 gutter</h4>
    <GridContainer>
      <Grid row gap={1}>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
      </Grid>
    </GridContainer>

    <h4>2 gutter</h4>
    <GridContainer>
      <Grid row gap={2}>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
      </Grid>
    </GridContainer>

    <h4>3 gutter</h4>
    <GridContainer>
      <Grid row gap={3}>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
      </Grid>
    </GridContainer>

    <h4>4 gutter</h4>
    <GridContainer>
      <Grid row gap={4}>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
      </Grid>
    </GridContainer>

    <h4>5 gutter</h4>
    <GridContainer>
      <Grid row gap={5}>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
      </Grid>
    </GridContainer>

    <h4>6 gutter</h4>
    <GridContainer>
      <Grid row gap={6}>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
        <Grid col={4}>{testContent}</Grid>
      </Grid>
    </GridContainer>
  </div>
);
export const DeviceBreakpoints = () => (
  <div>
    <GridContainer containerSize="card">
      <div style={exampleStyles}>card</div>
    </GridContainer>
    <GridContainer containerSize="card-lg">
      <div style={exampleStyles}>card-lg</div>
    </GridContainer>
    <GridContainer containerSize="mobile">
      <div style={exampleStyles}>mobile</div>
    </GridContainer>
    <GridContainer containerSize="mobile-lg">
      <div style={exampleStyles}>mobile-lg</div>
    </GridContainer>
    <GridContainer containerSize="tablet">
      <div style={exampleStyles}>tablet</div>
    </GridContainer>
    <GridContainer containerSize="tablet-lg">
      <div style={exampleStyles}>tablet-lg</div>
    </GridContainer>
    <GridContainer containerSize="desktop">
      <div style={exampleStyles}>desktop</div>
    </GridContainer>
    <GridContainer containerSize="desktop-lg">
      <div style={exampleStyles}>desktop-lg</div>
    </GridContainer>
    <GridContainer containerSize="widescreen">
      <div style={exampleStyles}>widescreen</div>
    </GridContainer>
  </div>
);
